import * as cdk from 'aws-cdk-lib';
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as s3deploy from 'aws-cdk-lib/aws-s3-deployment';
import * as cloudfront from 'aws-cdk-lib/aws-cloudfront';
import * as origins from 'aws-cdk-lib/aws-cloudfront-origins';
import * as route53 from 'aws-cdk-lib/aws-route53';
import * as route53targets from 'aws-cdk-lib/aws-route53-targets';
import * as acm from 'aws-cdk-lib/aws-certificatemanager';
import * as codepipeline from 'aws-cdk-lib/aws-codepipeline';
import * as codepipeline_actions from 'aws-cdk-lib/aws-codepipeline-actions';
import * as codebuild from 'aws-cdk-lib/aws-codebuild';
import * as iam from 'aws-cdk-lib/aws-iam';
import * as cloudwatch from 'aws-cdk-lib/aws-cloudwatch';
import * as cloudwatch_actions from 'aws-cdk-lib/aws-cloudwatch-actions';
import * as sns from 'aws-cdk-lib/aws-sns';
import * as targets from 'aws-cdk-lib/aws-events-targets';
import { Construct } from 'constructs';

export interface StaticSiteProps extends cdk.StackProps {
  /**
   * The domain name for the site (e.g. example.com)
   * @default - undefined
   */
  readonly domainName?: string;

  /**
   * The subdomain for the site (e.g. www)
   * @default - undefined
   */
  readonly siteSubDomain?: string;
  
  /**
   * The GitHub repository owner (e.g. your GitHub username)
   * @default - undefined
   */
  readonly repoOwner?: string;
  
  /**
   * The GitHub repository name
   * @default - undefined
   */
  readonly repoName?: string;
  
  /**
   * The GitHub personal access token
   * @default - undefined
   */
  readonly gitHubToken?: cdk.SecretValue;
}

export class StaticSiteStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props: StaticSiteProps) {
    super(scope, id, props);

    const siteDomain = props.domainName 
      ? (props.siteSubDomain ? `${props.siteSubDomain}.${props.domainName}` : props.domainName)
      : undefined;

    // Content bucket
    const siteBucket = new s3.Bucket(this, 'SiteBucket', {
      bucketName: siteDomain ? `${siteDomain}-static-site` : undefined,
      publicReadPolicy: true,
      websiteIndexDocument: 'index.html',
      websiteErrorDocument: 'error.html',
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      autoDeleteObjects: true,
    });

    new cdk.CfnOutput(this, 'Bucket', { value: siteBucket.bucketName });

    // CloudFront distribution
    let distribution: cloudfront.Distribution;
    
    if (siteDomain && props.repoOwner && props.repoName && props.gitHubToken) {
      // Get hosted zone
      const hostedZone = route53.HostedZone.fromLookup(this, 'Zone', {
        domainName: props.domainName!,
      });

      // Create certificate
      const certificate = new acm.DnsValidatedCertificate(this, 'SiteCertificate', {
        domainName: siteDomain,
        hostedZone,
        region: 'us-east-1', // CloudFront requires certificates in us-east-1
      });
      
      // Create distribution
      distribution = new cloudfront.Distribution(this, 'SiteDistribution', {
        defaultRootObject: 'index.html',
        domainNames: [siteDomain],
        certificate: certificate,
        defaultBehavior: {
          origin: new origins.S3Origin(siteBucket),
          compress: true,
          allowedMethods: cloudfront.AllowedMethods.ALLOW_GET_HEAD_OPTIONS,
          viewerProtocolPolicy: cloudfront.ViewerProtocolPolicy.REDIRECT_TO_HTTPS,
        },
      });

      // Add Route53 record
      new route53.ARecord(this, 'SiteAliasRecord', {
        recordName: siteDomain,
        target: route53.RecordTarget.fromAlias(new route53targets.CloudFrontTarget(distribution)),
        zone: hostedZone,
      });
      
      // Create deployment pipeline
      this.createPipeline(siteBucket, distribution, props);
      
      // Add monitoring
      this.addMonitoring(distribution, siteBucket);
    } else {
      // Create distribution without custom domain
      distribution = new cloudfront.Distribution(this, 'SiteDistribution', {
        defaultRootObject: 'index.html',
        defaultBehavior: {
          origin: new origins.S3Origin(siteBucket),
          compress: true,
          allowedMethods: cloudfront.AllowedMethods.ALLOW_GET_HEAD_OPTIONS,
          viewerProtocolPolicy: cloudfront.ViewerProtocolPolicy.REDIRECT_TO_HTTPS,
        },
      });
    }

    new cdk.CfnOutput(this, 'DistributionId', { value: distribution.distributionId });
    
    // Deploy site contents to bucket (for initial deployment)
    new s3deploy.BucketDeployment(this, 'DeployWithInvalidation', {
      sources: [s3deploy.Source.asset('../public')], // Adjust this path as needed
      destinationBucket: siteBucket,
      distribution,
      distributionPaths: ['/*'],
    });
  }
  
  private createPipeline(
    siteBucket: s3.Bucket,
    distribution: cloudfront.Distribution,
    props: StaticSiteProps
  ) {
    // Create build project
    const buildProject = new codebuild.PipelineProject(this, 'BuildProject', {
      buildSpec: codebuild.BuildSpec.fromObject({
        version: '0.2',
        phases: {
          install: {
            'runtime-versions': {
              nodejs: 14
            },
            commands: [
              'npm install -g aws-cdk',
              'npm install'
            ]
          },
          build: {
            commands: [
              'npm run build',
              'npm run test'
            ]
          }
        },
        artifacts: {
          'base-directory': 'dist',
          files: [
            '**/*'
          ]
        }
      }),
      environment: {
        buildImage: codebuild.LinuxBuildImage.STANDARD_5_0
      }
    });
    
    // Create pipeline
    const pipeline = new codepipeline.Pipeline(this, 'Pipeline', {
      pipelineName: 'MusingsStaticSitePipeline'
    });
    
    // Create source stage
    const sourceOutput = new codepipeline.Artifact();
    const sourceAction = new codepipeline_actions.GitHubSourceAction({
      actionName: 'GitHub_Source',
      owner: props.repoOwner!,
      repo: props.repoName!,
      branch: 'main',
      oauthToken: props.gitHubToken,
      output: sourceOutput
    });
    
    pipeline.addStage({
      stageName: 'Source',
      actions: [sourceAction],
    });
    
    // Create build stage
    const buildOutput = new codepipeline.Artifact();
    const buildAction = new codepipeline_actions.CodeBuildAction({
      actionName: 'Build',
      project: buildProject,
      input: sourceOutput,
      outputs: [buildOutput]
    });
    
    pipeline.addStage({
      stageName: 'Build',
      actions: [buildAction]
    });
    
    // Create deploy stage
    const deployAction = new codepipeline_actions.S3DeployAction({
      actionName: 'Deploy',
      bucket: siteBucket,
      input: buildOutput,
      runOrder: 1
    });
    
    // Create invalidate stage
    const invalidateRole = new iam.Role(this, 'InvalidateRole', {
      assumedBy: new iam.ServicePrincipal('codebuild.amazonaws.com')
    });
    
    distribution.grant(invalidateRole, 'cloudfront:CreateInvalidation');
    
    const invalidateProject = new codebuild.PipelineProject(this, 'InvalidateProject', {
      buildSpec: codebuild.BuildSpec.fromObject({
        version: '0.2',
        phases: {
          build: {
            commands: [
              'aws cloudfront create-invalidation --distribution-id ${CLOUDFRONT_ID} --paths "/*"'
            ]
          }
        }
      }),
      environmentVariables: {
        CLOUDFRONT_ID: {
          value: distribution.distributionId
        }
      },
      role: invalidateRole
    });
    
    const invalidateAction = new codepipeline_actions.CodeBuildAction({
      actionName: 'Invalidate',
      project: invalidateProject,
      input: sourceOutput,
      runOrder: 2
    });
    
    pipeline.addStage({
      stageName: 'Deploy',
      actions: [deployAction, invalidateAction]
    });
    
    new cdk.CfnOutput(this, 'PipelineName', { value: pipeline.pipelineName });
  }
  
  private addMonitoring(
    distribution: cloudfront.Distribution,
    siteBucket: s3.Bucket
  ) {
    // Create alarm topic
    const alarmTopic = new sns.Topic(this, 'AlarmTopic', {
      displayName: 'Static Site Monitoring Alerts',
      topicName: 'static-site-alarms'
    });
    
    // 4XX Error Rate Alarm
    const errorRateAlarm = new cloudwatch.Alarm(this, 'High4XXErrorRateAlarm', {
      metric: distribution.metric4XXErrorRate({
        period: cdk.Duration.minutes(5),
        statistic: 'Average'
      }),
      threshold: 50,
      evaluationPeriods: 2,
      comparisonOperator: cloudwatch.ComparisonOperator.GREATER_THAN_THRESHOLD,
      alarmDescription: 'Alarm when CloudFront 4XX error rate exceeds 50%',
      alarmName: 'CloudFrontHigh4XXErrorRate'
    });
    
    errorRateAlarm.addAlarmAction(new cloudwatch_actions.SnsAction(alarmTopic));
    
    // High Latency Alarm
    const latencyAlarm = new cloudwatch.Alarm(this, 'HighLatencyAlarm', {
      metric: distribution.metricOriginLatency({
        period: cdk.Duration.minutes(5),
        statistic: 'Average'
      }),
      threshold: 2000, // 2 seconds
      evaluationPeriods: 2,
      comparisonOperator: cloudwatch.ComparisonOperator.GREATER_THAN_THRESHOLD,
      alarmDescription: 'Alarm when CloudFront origin latency exceeds 2 seconds',
      alarmName: 'CloudFrontHighLatency'
    });
    
    latencyAlarm.addAlarmAction(new cloudwatch_actions.SnsAction(alarmTopic));
    
    // Low Bytes Downloaded Alarm (potential issue with content)
    const lowBytesAlarm = new cloudwatch.Alarm(this, 'LowBytesDownloadedAlarm', {
      metric: distribution.metricBytesDownloaded({
        period: cdk.Duration.minutes(5),
        statistic: 'Sum'
      }),
      threshold: 1000000, // 1MB
      evaluationPeriods: 2,
      comparisonOperator: cloudwatch.ComparisonOperator.LESS_THAN_THRESHOLD,
      alarmDescription: 'Alarm when bytes downloaded are less than 1MB',
      alarmName: 'CloudFrontLowBytesDownloaded'
    });
    
    lowBytesAlarm.addAlarmAction(new cloudwatch_actions.SnsAction(alarmTopic));
    
    // S3 Bucket Size Monitoring
    const bucketSizeMetric = new cloudwatch.Metric({
      namespace: 'AWS/S3',
      metricName: 'BucketSizeBytes',
      dimensionsMap: {
        BucketName: siteBucket.bucketName,
        StorageType: 'StandardStorage'
      },
      period: cdk.Duration.days(1),
      statistic: 'Average'
    });
    
    const bucketSizeAlarm = new cloudwatch.Alarm(this, 'BucketSizeAlarm', {
      metric: bucketSizeMetric,
      threshold: 1000000000, // 1GB
      evaluationPeriods: 1,
      comparisonOperator: cloudwatch.ComparisonOperator.GREATER_THAN_THRESHOLD,
      alarmDescription: 'Alarm when S3 bucket size exceeds 1GB',
      alarmName: 'S3BucketSizeExceeded'
    });
    
    bucketSizeAlarm.addAlarmAction(new cloudwatch_actions.SnsAction(alarmTopic));
    
    new cdk.CfnOutput(this, 'AlarmTopicArn', { value: alarmTopic.topicArn });
  }
}