package collector

/* 提供注册collector的方法registerCollector，同时要求要实现两个方法：
- signCollector，要求知道collector提供名字和指标名称
- exportMetric，将metric以map[string]any的形式返回
*/
