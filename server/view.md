# View

## /admin/edit
```
?[type]&[id]&[parentId]
```

> - 先将`type`, `id`, `parentId`进行`MustAtoi`
> - `type`为`1`:
>   - `id`为`[0-9][0-9]*`, 编辑对应文章
>   - `id`为其他, 在`parentId`对应目录下新建文章
> - `type`为其他, 显示`id`对应目录
