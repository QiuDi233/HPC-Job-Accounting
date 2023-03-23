# hpc-job-accounting

HPC 作业统计工具

## 基础功能

- 请在json文件中输入您的数据库相关信息

- 读取 hpc 作业日志文件

- 分析日志内容，抽取作业的信息

- 作业信息写入 mysql 数据库

## 作业信息

举例如下日志：

```
03/06/2023 01:38:30;E;124231.ml-hpc-master01;user=bilizhi group=users2 project=_pbs_project_default jobname=000_Main_C009_h3d queue=workq ctime=1678037519 qtime=1678037519 etime=1678037541 start=1678037544 exec_host=ml-rtx6000-ser002/1 exec_vnode=(ml-rtx6000-ser002:ncpus=1) Resource_List.higher=1 Resource_List.hwlic=6000 Resource_List.hwlimit=1 Resource_List.mpiprocs=1 Resource_List.ncpus=1 Resource_List.nodect=1 Resource_List.place=free Resource_List.select=1:ncpus=1:mpiprocs=1:pas_applications_enabled=GUI Resource_List.software=Hvtrans session=62548 end=1678037910 Exit_status=0 resources_used.cpupercent=135 resources_used.cput=00:06:49 resources_used.mem=534592kb resources_used.ncpus=1 resources_used.vmem=1362688kb resources_used.walltime=00:04:31 run_count=1
```

`03/06/2023 01:38:30;` log 时间，不用存

`E;` log 类型，Q-队列，S-开始，E-结束，D-删除，L-License。Q/S/E 类型的需要分析；D 和 L 类型的不用分析。

`124231.ml-hpc-master01;` 作业的 jobid = 124231 [ √ ]

`user=bilizhi` 用户名 [ √ ]

`jobname=000_Main_C009_h3d` 作业名 [ √ ]

`ctime=1678037519` ctime [ √ ]

`qtime=1678037519` qtime [ √ ]

`etime=1678037541` etime [ √ ]

`start=167803754` 开始运行时刻 start_time [ √ ]

`Resource_List.ncpus=1` 消耗 cpu 数量 [ √ ]

`Resource_List.nodect=1` 消耗计算节点数量 [ √ ]

`Resource_List.software=Hvtrans` 软件名 software [ √ ]

`end=1678037910` 结束时刻 end_time [ √ ]

`Exit_status=0` 程序退出码 [ √ ]


上面日志里没有，其他日志有的：

`Resource_List.select=3:ncpus=36:mpiprocs=36:pas_applications_enabled=Dyna:platform=G6254-36c-384G` 

其中 platform=G6254-36c-384G [ √ ]

## 涉及知识点

- git 使用

- 读文件

- 字符串处理

- 数据库
