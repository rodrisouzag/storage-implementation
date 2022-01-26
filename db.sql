create table products(
    id int not null primary key auto_increment,
    name text not null,
    category text not null,
    count int not null,
    price float not null
);