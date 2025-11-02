import * as path from 'path';
import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import { RestApi, LambdaIntegration } from 'aws-cdk-lib/aws-apigateway';
import * as iam from 'aws-cdk-lib/aws-iam';

export class UserManagerStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // Secrets from SSM Parameter Store
    const jwtParamName = '/user-manager/app/jwt-secret';
    const apiKeyParamName = '/user-manager/app/api-key';

    // Lambda
    const appLambda = new lambda.Function(this, 'UserManagerHandler', {
      runtime: lambda.Runtime.PROVIDED_AL2023,
      handler: 'bootstrap',
      code: lambda.Code.fromAsset(path.join(__dirname, '../../lambdas')),
      architecture: lambda.Architecture.ARM_64,
      memorySize: 512,
      timeout: cdk.Duration.seconds(15),
      environment: {
        APP_ENV: 'production',
        APP_PORT: '8080',
        APP_BASE_URL: 'http://user-manager.com',
        APP_JWT_SECRET_PARAM: jwtParamName,
        APP_API_KEY_PARAM: apiKeyParamName,
        MAIL_FROM_EMAIL: 'dangeschichte@gmail.com',
        SES_REGION: 'eu-central-1'
      },
    });
    // Role policy for secret params:
    appLambda.addToRolePolicy(new iam.PolicyStatement({
      actions: ['ssm:GetParameter'],
      resources: [
        `arn:aws:ssm:${this.region}:${this.account}:parameter${jwtParamName}`,
        `arn:aws:ssm:${this.region}:${this.account}:parameter${apiKeyParamName}`,
      ],
    }));

    // DynamoDB
    const usersTable = new dynamodb.TableV2(this, 'UserManagerUsersTable', {
      tableName: 'users',
      partitionKey: { name: 'id', type: dynamodb.AttributeType.STRING },
      billing: dynamodb.Billing.onDemand(),
      globalSecondaryIndexes: [
        {
          indexName: 'email-index',
          partitionKey: { name: 'email', type: dynamodb.AttributeType.STRING },
          projectionType: dynamodb.ProjectionType.ALL,
        },
      ],
    });
    usersTable.grantReadWriteData(appLambda);

    // SES 
    appLambda.addToRolePolicy(new iam.PolicyStatement({
      actions: ['ses:SendEmail', 'ses:SendRawEmail'],
      resources: ['*'],
    }));

    // API Gateway
    const api = new RestApi(this, 'UserManagerApi', {
      defaultCorsPreflightOptions: {
        allowOrigins: ['*'],
        allowMethods: ['OPTIONS', 'GET', 'POST', 'PUT', 'DELETE'],
        allowHeaders: ['Content-Type', 'Authorization'],
        allowCredentials: false,
      },
    });
    api.root.addProxy({
      defaultIntegration: new LambdaIntegration(appLambda, { proxy: true }),
      anyMethod: true,
    });

    // Sets APP_BASE_URL after deployment
    appLambda.addEnvironment('APP_STAGE', 'prod');
  }
}
