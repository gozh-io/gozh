# gozh
获取项目
git clone https://github.com/gozh-io/gozh.git


说明1:gozh中文社区后端

说明2:每次修改完代码后,需要先在本机上执行godep save,保存gozh引用的依赖到vendor目录里

说明3:修改代码,修改配置,都需要重新编译,发布image, 编译之前需要安装make工具


1.编译gozh image
make

2.发布到 docker hub
make publish
有时候会提示超时,重新执行发布
执行发布时,修改Makefile中REGISTRY_USER,改成自己的账号

3.测试是否发布成功
make publish-test

4.删除本地和已发布到docker hub上的iamge
make rmi

5.编译,发布,发布测试
make all


gozh后端api暴露了下面内容
VOLUME /myapp 
EXPOSE 80

运行一个docker container
docker run --rm -d -p 80:80 blade2iron/gozh
Makefile中 REGISTRY_USER 默认是blade2iron,如果修改 REGISTRY_USER , 请把blade2iron改成对应的值

查看运行docker container
docker ps

ssh到一个已运行container上
docker exec -it e0067e87bdfa sh
其中e0067e87bdfa 是container id

