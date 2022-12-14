create table auth_account(
    id bigint(20) unsigned not null primary key auto_increment comment '主键',
    created_at datetime not null default now() comment '创建时间',
    created_by varchar(64) default '' comment '创建人',
    updated_at datetime not null default now() comment '修改时间',
    updated_by varchar(64) default '' comment '修改人',
    is_deleted tinyint(3) default 0 comment '是否删除 0未删除 1已删除',
    
    username varchar(128) default '' comment '用户名',
    `password` varchar(128) default '' comment '密码',
    role tinyint(3) unsigned DEFAULT 0 comment '角色 0 普通用户、1 管理员',

    state tinyint(3) unsigned DEFAULT 1 comment '状态 0 为禁用、1 为启用'
)engine=innodb charset=utf8mb4 comment '用户信息表';