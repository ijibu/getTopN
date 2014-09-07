## 海量日志数据，提取出访问次数最多的N个IP
既然是海量数据处理，那么可想而知，给我们的数据那就一定是海量的。针对这个数据的海量，我们如何着手呢?对的，无非就是分而治之/hash映射 + hash统计 + 堆/快速/归并排序，说白了，就是先映射，而后统计，最后排序：
 - 代码约束
- ip.log为已知的，即为一个大文件，每行存储的是一个IP。请去[百度网盘](http://pan.baidu.com/s/1kTDn8f9)下载，提取码:tcih
- log目录，下面放的是对ip.log进行分割后的小文件。

有何问题或补充意见，请联系[ijibu](http://weibo.com/ijibu)  。
