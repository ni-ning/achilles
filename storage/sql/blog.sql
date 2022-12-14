create table blog_article(
    id bigint(20) unsigned not null primary key auto_increment comment '主键',
    created_at datetime not null default now() comment '创建时间',
    created_by varchar(64) default '' comment '创建人',
    updated_at datetime not null default now() comment '修改时间',
    updated_by varchar(64) default '' comment '修改人',
    -- deleted_at datetime null comment '删除时间',
    is_deleted tinyint(3) default 0 comment '是否删除 0未删除 1已删除',
    
    title varchar(100) default '' comment '文章标题',
    `desc` varchar(255) default '' comment '文章简述',
    cover_image_url varchar(255) default '' comment '封面图片地址',
    content longtext comment '文章内容',
    state tinyint(3) unsigned DEFAULT 1 comment '状态 0 为禁用、1 为启用'
)engine=innodb charset=utf8mb4 comment '文章表';


create table blog_tag(
    id bigint(20) unsigned not null primary key auto_increment comment '主键',
    created_at datetime not null default now() comment '创建时间',
    created_by varchar(64) default '' comment '创建人',
    updated_at datetime not null default now() comment '修改时间',
    updated_by varchar(64) default '' comment '修改人',
    -- deleted_at datetime null comment '删除时间',
    is_deleted tinyint(3) default 0 comment '是否删除 0未删除 1已删除',

    name varchar(100) default '' comment '标签名称',
    state tinyint(3) unsigned DEFAULT 1 comment '状态 0 为禁用、1 为启用'
)engine=innodb charset=utf8mb4 comment '标签表';


create table blog_article_tag(
    id bigint(20) unsigned not null primary key auto_increment comment '主键',
    created_at datetime not null default now() comment '创建时间',
    created_by varchar(64) default '' comment '创建人',
    updated_at datetime not null default now() comment '修改时间',
    updated_by varchar(64) default '' comment '修改人',
    -- deleted_at datetime null comment '删除时间',
    is_deleted tinyint(3) default 0 comment '是否删除 0未删除 1已删除',
    
    article_id bigint(20) unsigned not null comment '文章 ID',
    tag_id bigint(20) unsigned not null comment '标签 ID'
)engine=innodb charset=utf8mb4 comment '文章标签关联表';