# 接口文档

## `/api/admin/new`

```toml
{
  data string
  type int (=0|1)
  dirId int
}
```

```toml
{
  res string (="ok"|"err")
}
```

> `dirId`默认为0
>
> `type`=0, `data`为目录名  
> `type`=1, `data`为markdown
> 
> 上传文件/图片使用`newResource`

## `/api/admin/newResource`

```toml
Content-Type: application/octet-stream
-------------------------------
       DATA: {FileNameLen} {FileName} {IsTmp} {NonTmpProps} {FileData}
FileNameLen: <[4]uint32>
   FileName: <[FileNameLen]string(UTF-8)>
      IsTmp: <[1]bool>
NonTmpProps: IsTmp ? null : {DirId}
      DirId: <[4]uint32>
   FileData: <[*]binary>
```

```toml
{
  res string (="ok"|"err")
  fileLoc string
}
```

## `/api/admin/edit`

```toml
{
  data string
  type int (=0|1|2)
  id int
}
```

```toml
{
  res string (="ok"|"err")
}
```

> type=0, data为目录名  
> type=1, data为markdown  
> type=2, data为文件名

## `/api/admin/remove`

```toml
{
  list: []{
    type int
    id int
  }
}
```

```toml
{
  res string (="ok"|"err")
}
```

## `/api/admin/move`

```toml
{
  list: []{
    type int
    id int
  }
  dirId int
}
```

```toml
{
  res string (="ok"|"err")
}
```

> `dirId`默认为`0`
