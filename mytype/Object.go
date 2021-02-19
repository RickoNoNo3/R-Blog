// map不能直接等号赋值, 要用指针, 因此封装了get/set函数
package mytype

// 将多个对象类型合为一个键值对组所用的类型
type ObjectList map[string]*Object

// 仿JS对象类型 (map[string]interface{}的便利化封装)
// 实现多层嵌套键值对
type Object struct {
	Val  interface{}
	list *ObjectList
}

// 判断一个Object是否拥有某一key
func (obj *Object) Has(key string) bool {
	return (*obj.list)[key] != nil
}

// 获取一个Object中一个key对应的子Object, 未找到key就会panic
func (obj *Object) Get(key string) (val *Object) {
	if !obj.Has(key) {
		panic("No Such Key!")
	}
	return (*obj.list)[key]
}

// 给一个Object中的一个key设置对应的子Object, 可同时用于创建和更新
func (obj *Object) Set(key string, val *Object) (valRes *Object) {
	(*obj.list)[key] = val
	return obj.Get(key)
}

// 静态化. 将Object转为map[string]interface{}结构, 便于在其他地方使用(如构造json/模板传参).
// 注意有子元素的会看作Group节点, 否则才看作Value节点, 因此有子元素的元素的Val不会被导出.
func (obj *Object) Staticize() (res map[string]interface{}) {
	res = make(map[string]interface{})
	for k, v := range *obj.list {
		if len(*v.list) == 0 { // 不是Group
			res[k] = v.Val
		} else { // 是Group
			res[k] = v.Staticize()
		}
	}
	return
}

// 新建一个普通Object对象, 注意有children的对象的Val在静态化时无效
func NewObject(value interface{}, children ObjectList) *Object {
	return &Object{
		Val:  value,
		list: &children,
	}
}

// 新建一个Group对象, 有子元素, 无自身值
func NewGroup(children ObjectList) *Object {
	return NewObject(nil, children)
}

// 新建一个Value对象, 无子元素, 有自身值
func NewValue(value interface{}) *Object {
	return NewObject(value, ObjectList{})
}
