CREATE TABLE `article` (
	`id`	INTEGER NOT NULL,
	`title`	TEXT NOT NULL DEFAULT "",
	`markdown`	TEXT NOT NULL DEFAULT "",
	`html`	TEXT NOT NULL DEFAULT "",
	`contentsJ`	TEXT NOT NULL DEFAULT "",
	`commentJ`	TEXT NOT NULL DEFAULT "",
	`drewT`	TEXT NOT NULL DEFAULT "",
	`tagJ`	TEXT NOT NULL DEFAULT "",
	PRIMARY KEY(`id` AUTOINCREMENT)
);
CREATE TABLE `dir` (
	`id`	INTEGER NOT NULL,
	`name`	TEXT NOT NULL DEFAULT "",
	PRIMARY KEY(`id` AUTOINCREMENT)
);
CREATE TABLE `layer` (
	`id`	INTEGER NOT NULL,
	`type`	INTEGER NOT NULL,
	`dirId`	INTEGER NOT NULL DEFAULT 0,
	`createdT`	TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`modifiedT`	TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY(`id`,`type`)
);
CREATE TRIGGER article_insert_tri AFTER INSERT ON article
BEGIN
  INSERT INTO layer(id, type) VALUES(new.id, 0);
  -- 这个能递归
  UPDATE layer SET modifiedT=CURRENT_TIMESTAMP WHERE id=new.id AND type=0;
END;
CREATE TRIGGER dir_insert_tri AFTER INSERT ON dir
BEGIN
  INSERT INTO layer(id, type) VALUES(new.id, 1);
  -- 这个能递归
  UPDATE layer SET modifiedT=CURRENT_TIMESTAMP WHERE id=new.id AND type=1;
END;
CREATE TRIGGER article_update_tri_html AFTER UPDATE OF html ON article
BEGIN
  UPDATE article SET drewT=CURRENT_TIMESTAMP WHERE id=new.id;
END;
CREATE TRIGGER article_update_tri_markdown AFTER UPDATE OF markdown ON article
BEGIN
  UPDATE layer SET modifiedT=CURRENT_TIMESTAMP WHERE id=new.id AND type=0;
END;
CREATE TRIGGER layer_update_tri_modifiedT AFTER UPDATE OF modifiedT ON layer
BEGIN
  -- 这个能递归
  UPDATE layer SET modifiedT=CURRENT_TIMESTAMP WHERE id=new.dirId AND type=1;
END;
CREATE TRIGGER article_delete_tri AFTER DELETE ON article
BEGIN
  DELETE FROM layer WHERE id=old.id AND type=0;
END;
CREATE TRIGGER dir_delete_tri AFTER DELETE ON dir
BEGIN
  -- 这个能递归
  DELETE FROM article WHERE id IN (SELECT id FROM layer WHERE dirId=old.id AND type=0);
  DELETE FROM dir WHERE id IN (SELECT id FROM layer WHERE dirId=old.id AND type=1);
  -- 删除自身
  DELETE FROM layer WHERE id=old.id AND type=1;
END;
CREATE TRIGGER layer_update_tri_dirId AFTER UPDATE OF dirId ON layer
BEGIN
  UPDATE layer SET modifiedT=CURRENT_TIMESTAMP WHERE id=new.dirId AND type=1;
END;
CREATE VIEW layer_read as
select dirId, layer.id, type, name as `text`, createdT, modifiedT
  from layer inner join dir on layer.id = dir.id and type=1
union
select dirId, layer.id, type, title as `text`, createdT, modifiedT
  from layer inner join article on layer.id = article.id and type=0 order by dirId asc, type desc, modifiedT desc;
