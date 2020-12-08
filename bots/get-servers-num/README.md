# サーバ台数をチャンネルに取得 & 表示するだけのボット

## 使い方

```
@<bot_name> get-servers-num
現在のサーバー台数は... 0台です！！！１
```

## デプロイ方法

ボット起動用 EC2 インスタンスに `/home/ec2-user/.infra-bot.yaml` というファイル名で Slack app の token を記す

```
token: xoxb-****-****-****
```

ボット起動用 EC2 インスタンスにて書きを実行

```shell
aws ecr get-login --region ap-northeast-1 --no-include-email
# 出力されたメッセージの 'docker login -u ...' を入力
sudo docker login -u AWS -p ....
sudo docker run -d "/home/ec2-user/.infra-bot.yaml:/root/.infra-bot.yaml" xxxxxx.dkr.ecr.ap-northeast-1.amazonaws.com/get-servers:<tag>
```
