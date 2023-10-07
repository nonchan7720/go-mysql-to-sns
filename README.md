# go-mysql-to-sns

MySQLのBinlogを監視し

- 新しいレコードの挿入（Create）
- 既存のレコードの更新（Update）
- レコードの削除（Delete）

の操作を検出し、それらの変更をAWS SNSトピックに送信します。
