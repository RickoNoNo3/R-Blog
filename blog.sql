PRAGMA foreign_keys= OFF;
BEGIN TRANSACTION;

CREATE TABLE `dir` (
  `id`    INTEGER NOT NULL,
  `title` TEXT    NOT NULL,
  PRIMARY KEY (`id` AUTOINCREMENT)
);

CREATE TABLE `file` (
  `id`    INTEGER NOT NULL,
  `title` TEXT    NOT NULL,
  PRIMARY KEY (`id` AUTOINCREMENT)
);

CREATE TABLE `article` (
  `id`       INTEGER NOT NULL,
  `title`    TEXT    NOT NULL,
  `markdown` TEXT    NOT NULL,
  `tags`     TEXT    NOT NULL DEFAULT '',
  `voted`    INTEGER NOT NULL DEFAULT 0,
  `visited`  INTEGER NOT NULL DEFAULT 0,
  PRIMARY KEY (`id` AUTOINCREMENT)
);

CREATE TABLE `layer` (
  `id`        INTEGER NOT NULL,
  `type`      INTEGER NOT NULL,
  `dirId`     INTEGER NOT NULL DEFAULT 0,
  `createdT`  TEXT    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `modifiedT` TEXT    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`, `type`)
);

CREATE VIEW layer_read AS
  SELECT dirId, layer.id, type, title, createdT, modifiedT
  FROM layer
         INNER JOIN dir ON layer.id = dir.id AND type = 0
  UNION
  SELECT dirId, layer.id, type, title, createdT, modifiedT
  FROM layer
         INNER JOIN article ON layer.id = article.id AND type = 1
  UNION
  SELECT dirId, layer.id, type, title, createdT, modifiedT
  FROM layer
         INNER JOIN file ON layer.id = file.id AND type = 2
  ORDER BY dirId ASC, type ASC, modifiedT DESC, createdT DESC;

INSERT INTO article
VALUES
  (0, '建博客的一些心得', REPLACE(REPLACE('# 建博客的一些心得\r\012\r\012欢迎你来我的博客！\r\012\r\012早在刚上大学的时候就听学长说过，如果我有一个自己的博客，而且经常在里面写东西的话，对自己会很有帮助。但我又不是很想在cnblog或者csdn或者github种种平台上搞个博客就那么随便用，最终还是决定自己搭一个。\r\012\r\012然而，由于人类的本质是鸽子（不是）的一部分原因，搭建这个博客的时候我已经大三了。毕竟也已经大三了，前端的部分学的差不多了，后端的Golang也算会用，所以搭建博客不是一个特别棘手的问题，更像是把自己的所学（包括软件工程的内容）逐步尝试应用到实际中来，作为一个技术经验的积累。\r\012\r\012当然写文章也是很重要的。搭好了博客总得写点什么东西，这样才叫充实。说到写文章，我在这个博客里写的东西应该会比较杂，而且这个博客又是在不成熟的技术条件下搭建的，分类功能难免有些缺陷。自己的专业水平又不是和大犇们一样强，写的东西可能也没什么水平。我会尽量克服这些问题，也请见谅。\r\012\r\012玩的开心~！\r\012\r\0122019.9.9\r\012\r\012', '\r', char(13)), '\012', char(10)), '', 0, 0);

DELETE FROM sqlite_sequence;

COMMIT;
