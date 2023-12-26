create table if not exists sensor_data (
    `id1` varchar(255) not null,
    `id2` bigint(20) unsigned not null,
    `type` varchar(255) not null,
    `timestamp` timestamp not null,
    `value` float not null,
    constraint force_upper_case check(binary id1 = upper(id1))
);

alter table sensor_data add index(`timestamp`);
alter table sensor_data add index(`id1`,`id2`);
alter table sensor_data add index(`id1`,`id2`,`timestamp`);