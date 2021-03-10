package chain

/**迭代器（iterator）是一种检查容器内元素并遍历元素的数据类型。
*定义迭代器接口标准，通过总结分析，迭代器有两个功能：
	1、判断容器中是否还有数据
	2、从容器中取出一个数据
 */
type Iterator interface {
	HasNext()bool //判断是否还有下一个数据
	Next()Block //返回下一个区块
}