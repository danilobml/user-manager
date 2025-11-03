import { App } from 'aws-cdk-lib';
import { Template, Match } from 'aws-cdk-lib/assertions';
import { UserManagerStack } from '../lib/user-manager-stack';

function synth() {
    const app = new App();
    const stack = new UserManagerStack(app, 'UserManagerStackTest', {
        env: { account: '111111111111', region: 'eu-central-1' },
    });
    return Template.fromStack(stack);
}

describe('UserManagerStack', () => {
    test('Lambda is created with expected properties', () => {
        const template = synth();
        template.hasResourceProperties('AWS::Lambda::Function', {
            Runtime: 'provided.al2023',
            Handler: 'bootstrap',
            Environment: {
                Variables: Match.objectLike({
                    APP_JWT_SECRET_PARAM: '/user-manager/app/jwt-secret',
                }),
            },
        });
    });

    test('API Gateway method integrates with Lambda', () => {
        const template = synth();
        template.hasResourceProperties('AWS::ApiGateway::Method', {
            HttpMethod: 'ANY',
            Integration: Match.objectLike({
                Type: 'AWS_PROXY',
                IntegrationHttpMethod: 'POST',
                Uri: Match.objectLike({ 'Fn::Join': Match.anyValue() }),
            }),
        });
    });

    test('API Gateway can invoke the Lambda', () => {
        const template = synth();
        template.hasResourceProperties('AWS::Lambda::Permission', {
            Action: 'lambda:InvokeFunction',
            Principal: 'apigateway.amazonaws.com',
        });
    });

    test('DynamoDB table "users" with email-index', () => {
        const template = synth();

        template.hasResourceProperties('AWS::DynamoDB::GlobalTable', {
            GlobalSecondaryIndexes: Match.arrayWith([
                Match.objectLike({ IndexName: 'email-index' }),
            ]),
        })
    });
});
