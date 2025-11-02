#!/usr/bin/env node
import * as cdk from 'aws-cdk-lib/core';
import { UserManagerStack } from '../lib/user-manager-stack';

const app = new cdk.App();
new UserManagerStack(app, 'UserManagerStack', {});
