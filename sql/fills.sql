INSERT INTO country(country_name) VALUES
    ('Дания'),
    ('Швеция'),
    ('Великобритания'),
    ('США'),
    ('Канада'),
    ('Франция');


INSERT INTO genre(genre_name) VALUES
    ('Боевик'),
    ('Драма'),
    ('Комедия'),
    ('Хоррор'),
    ('Мелодрама'),
    ('Триллер'),
    ('Вестерн'),
    ('Документальный'),
    ('Фентези'),
    ('Приключения');


INSERT INTO profile(first_name, surname, email, password, picture_url, gender, register_date) VALUES
    ('Иван', 'Иванов', 'ivan@mail.ru', '$2a$10$XFUgsxdBN.UiILtfITr/4urH1WIQWBMkqvLnAgfiYpZNguvySCBAq', '/static/media/img/users/base.jpg', 'male', current_timestamp),
    -- 1234abcd
    ('Алексей','Самойлов','lexa@mail.ru','$2a$10$c5Ku7dxUrlqqxCox6iRMLuKH8m7ohsBFv8OXIcEStH40/HqHuDAUO','/static/media/img/users/base.jpg','male',current_timestamp),
    -- 123ab321bc
    ('Катя','Кириллова','kate@mail.ru','$2a$10$ioIHXpmM2l0PMmd/4gJs..9yqt0O9kzklrIZoSTnZRnm2mmC8AMAi','/static/media/img/users/base.jpg','female',current_timestamp),
    -- ab1423bc
    ('Настя','Куликова','kulikovAnast@mail.ru','$2a$10$gATPQVuopB77CwgXWfR1behANZ9KI.nH2vPDIJTvN1yuIdEcQtcPa','/static/media/img/users/base.jpg','female',current_timestamp),
    --kkk111kkk
    ('Максим','Дудник','maksongold@mail.ru','2a$10$T62vmMhEeRKFLF82O6YhTOx8VmHKsYWximYpF3y7qLbTV4tPIFBEa','/static/media/img/users/base.jpg','male',current_timestamp),
    --maks2000
    ('Борис','Кожуро','borisK@mail.ru','2a$10$J2qOuCh9e9h4OE7DtQ7lZ.969oLF0B1N3HLWoq/tQg/zx38qFwd1y','/static/media/img/users/base.jpg','male',current_timestamp),
    --bk1234ev
    ('Анна','Морозова','morozAN@gmail.com','2a$10$t.QpAwfS.nj3m/v59z17FOsCbJzBdL0XipAVK4ecEFp9r4EUOVCYO','/static/media/img/users/base.jpg','female',current_timestamp),
    --moroz1523
    ('Игорь','Николаевич','IgorNikol@yandex.ru','$2a$10$sHl.NJcJTSeGqvJbmRXUQ.iOez7Ah4U4kFNKbAAB1YcqeMrbluzlq','/static/media/img/users/base.jpg','male',current_timestamp);
    -- kolya1967
    
    
INSERT INTO collection(author_id, collection_name, description, creation_time, picture_url) VALUES
    (1, 'Для ценителей Хогвартса','',current_timestamp,'server/images/collections1.png'),
    (1, 'Про настоящую любовь','',current_timestamp,'server/images/collections2.png'),
    (1, 'Аферы века','',current_timestamp,'server/images/collections3.png'),
    (2, 'Про Вторую Мировую','',current_timestamp,'server/images/collections4.png'),
    (2, 'Осеннее настроение','',current_timestamp,'server/images/collections5.png'),
    (2, 'Летняя атмосфер','',current_timestamp,'server/images/collections6.png'),
    (6, 'Про дружбу','',current_timestamp,'server/images/collections7.png'),
    (6, 'Романтические фильмы','',current_timestamp,'server/images/collections8.png'),
    (6, 'Джунгли зовут','',current_timestamp,'server/images/collections9.png'),
    (6, 'Фантастические фильмы','',current_timestamp,'server/images/collections10.png'),
    (3, 'Про петлю времени','',current_timestamp,'server/images/collections11.png'),
    (1, 'Классика на века','',current_timestamp,'server/images/collections1.png');


INSERT INTO person VALUES 
    (DEFAULT, 'Thomas Vinterberg', 'Томас Винтерберг', '', 'Режиссер, Сценарист', 183, 52, '1969-05-19', NULL ,'Копенгаген, Дания','','Муж.','Хелена Рейнгор Нойманн', NULL),
    (DEFAULT, 'Christian Ditter', 'Кристиан Диттер', '', 'Режиссер, Сценарист, Продюсер, Монтажер', 0, 44, '1977-01-01',NULL,'Гисен, Германия','','Муж.','',NULL),
    (DEFAULT, '', 'Сесилия Ахерн', '', '', 180, 20, NULL,NULL,'','','male','',NULL),
    (DEFAULT, '', 'Мадс Миккельсен', '', '', 180, 20, NULL,NULL,'','','male','',NULL),
    (DEFAULT, '', 'Томас Бо Ларсен', '', '', 180, 20, NULL,NULL,'','','male','',NULL),
    (DEFAULT, '', 'Ларс Ранте', '', '', 180, 20, NULL,NULL,'','','male','',NULL),
    (DEFAULT, '', 'Мария Бонневи', '', '', 180, 20, NULL,NULL,'','','male','',NULL),
    (DEFAULT, '', 'Лили Коллинз', '', '', 180, 20, NULL,NULL,'','','male','',NULL),
    (DEFAULT, '', 'Сэм Клафлин', '', '', 180, 20, NULL,NULL,'','','male','',NULL),
    (DEFAULT, '', 'Кристиан Кук', '', '', 180, 20, NULL,NULL,'','','male','',NULL),
    --10 up
    (DEFAULT, '', 'Джо Берлингер', '', '', 180, 20, NULL,NULL,'','','male','',NULL),
    (DEFAULT, '', 'Майкл Верви', '', '', 180, 20, NULL,NULL,'','','male','',NULL),
    (DEFAULT, '', 'Зак Эфрон', '', '', 180, 20, NULL,NULL,'','','male','',NULL),
    (DEFAULT, '', 'Том Шенклэнд', '', '', 180, 20, NULL,NULL,'','','male','',NULL),
    (DEFAULT, '', 'Виктор Гюго', '', '', 180, 20, NULL,NULL,'','','male','',NULL),
    (DEFAULT, '', 'Доминик Уэст', '', '', 180, 20, NULL,NULL,'','','male','',NULL),
    (DEFAULT, '', 'Тарсем Сингх', '', '', 180, 20, NULL,NULL,'','','male','',NULL),
    (DEFAULT, '', 'Марк Клейн', '', '', 180, 20, NULL,NULL,'','','male','',NULL),
    (DEFAULT, '', 'Джулия Робертс', '', '', 180, 20, NULL,NULL,'','','male','',NULL),
    (DEFAULT, '', 'Нэйтан Лейн', '', '', 180, 20, NULL,NULL,'','','male','',NULL),
    --20 up
    (DEFAULT, '', 'Крис Коламбус', '', '', 180, 20, NULL,NULL,'','','male','',NULL),
    (DEFAULT, '', 'Дж.К. Роулинг', '', '', 180, 20, NULL,NULL,'','','male','',NULL),
    (DEFAULT, '', 'Дэниэл Рэдклифф', '', '', 180, 20, NULL,NULL,'','','male','',NULL),
    (DEFAULT, '', 'Руперт Гринтн', '', '', 180, 20, NULL,NULL,'','','male','',NULL),
    (DEFAULT, '', 'Эмма Уотсон', '', '', 180, 20, NULL,NULL,'','','male','',NULL),
    (DEFAULT, '', 'Альфонсо Куарон', '', '', 180, 20, NULL,NULL,'','','male','',NULL);
    --26 up


INSERT INTO film VALUES
    (DEFAULT, 'Еще по одной', 'Druk', 0, 'В ресторане собираются учитель истории, психологии, музыки и физрук, чтобы отметить 40-летие одного из них. И решают проверить научную теорию о том, что c самого рождения человек страдает от нехватки алкоголя в крови, а чтобы стать по-настоящему счастливым, нужно быть немного нетрезвым. Друзья договариваются наблюдать, как возлияния скажутся на их работе и личной жизни, и устанавливают правила: не пить вечером и по выходным. Казалось бы, что может пойти не так?', 'server/images/one-more-drink.webp', 'trailer', 12000000.0, 2020, 132, 1, 1, 'фильм'),
    (DEFAULT, 'С любовью, Рози', 'Love, Rosie', 0, 'Рози и Алекс были лучшими друзьями с детства, и теперь, по окончании школы, собираются вместе продолжить учёбу в университете.\n\nОднако в их судьбах происходит резкий поворот, когда после ночи со звездой школы Рози узнаёт, что у неё будет ребенок. Невзирая на то, что обстоятельства и жизненные ситуации разлучают героев, они и спустя годы продолжают помнить друг о друге и о том чувстве, что соединило их в юности…', 'server/images/love-rosie.webp', 'trailer', 25384687.0, 2014, 132, 3, 2, 'фильм'),
    (DEFAULT, 'Красивый, плохой, злой', 'xtremely Wicked, Shockingly Evil and Vile', 0, 'Тэд – неотразимый красавец, устоять перед таким невозможно. Очаровательная Лиз и не устояла. Свидание на одну ночь переросло в долгие замечательные отношения, о которых большинство девушек только и мечтает. Внезапно Тэда арестовали, предъявив ему чудовищные обвинения. Но как поверить в то, что этот нежный умный элегантный мужчина мог... насиловать, убивать и расчленять несчастных женщин?! Ее родной и милый Тэд. Тэд Банди.', 'server/images/extremely-wicked.webp', 'trailer', 9816572.0, 2018, 110, 12, 11, 'фильм'),
    (DEFAULT, 'Отверженные', 'Les Misérables', 0, 'Жан Вальжан – осужденный за кражу хлеба беглый каторжник, на протяжении долгих лет скрывается от инспектора Жавера. Жавер, в свою очередь, свято верит в справедливость закона и считает делом чести поимку беглого вора. Параллельно с их противостоянием развивается трагичная история Фантины и ее незаконнорожденной дочери Козетты. После смерти Фантины Жан Вальжан, считающий себя в ответе за ее судьбу, забирает Козетту и заботится о ней как о родной дочери.', 'server/images/les-m.webp', 'trailer', 9816572.0, 2019, 360, 15, 14, 'фильм'),
    (DEFAULT, 'Белоснежка: Месть гномов', 'Mirror Mirror', 0, 'Злая Королева, мечтающая выйти замуж за красивого и богатого Принца, выдворяет из дворца Белоснежку и берет власть в свои руки. Но милая девушка не погибла в темном дремучем лесу, а связалась с бандой гномов-разбойников. Вместе они отомстят Злодейке!', 'server/images/mirror-mirror.webp', 'trailer', 85000000.0, 2012, 106, 18, 17, 'фильм'),
    (DEFAULT, 'Гарри Поттер и философский камень', 'Harry Potter and the Sorcerer''s Stone', 0, 'Жизнь десятилетнего Гарри Поттера нельзя назвать сладкой: родители умерли, едва ему исполнился год, а от дяди и тёти, взявших сироту на воспитание, достаются лишь тычки да подзатыльники. Но в одиннадцатый день рождения Гарри всё меняется. Странный гость, неожиданно появившийся на пороге, приносит письмо, из которого мальчик узнаёт, что на самом деле он - волшебник и зачислен в школу магии под названием Хогвартс. А уже через пару недель Гарри будет мчаться в поезде Хогвартс-экспресс навстречу новой жизни, где его ждут невероятные приключения, верные друзья и самое главное — ключ к разгадке тайны смерти его родителей.', 'server/images/harry1.webp', 'trailer', 974755371.0, 2001, 152, 22, 21, 'фильм'),
    (DEFAULT, 'Гарри Поттер и Тайная комната', 'Harry Potter and the Chamber of Secrets', 0, 'Гарри Поттер переходит на второй курс Школы чародейства и волшебства Хогвартс. Эльф Добби предупреждает Гарри об опасности, которая поджидает его там, и просит больше не возвращаться в школу. Юный волшебник не следует совету эльфа и становится свидетелем таинственных событий, разворачивающихся в Хогвартсе. Вскоре Гарри и его друзья узнают о существовании Тайной Комнаты и сталкиваются с новыми приключениями, пытаясь победить темные силы.', 'server/images/harry2.webp', 'trailer', 878969634.0, 2002, 161, 22, 21, 'фильм'),
    (DEFAULT, 'Гарри Поттер и узник Азкабана', 'Harry Potter and the Prisoner of Azkaban', 0, 'третьей части истории о юном волшебнике полюбившиеся всем герои — Гарри Поттер,Рон и Гермиона — возвращаются уже на третий курс школы чародейства и волшебства Хогвартс. На этот раз они должны раскрыть тайну узника, сбежавшего из зловещей тюрьмы Азкабан, чье пребывание на воле создает для Гарри смертельную опасность...', 'server/images/harry3.webp', 'trailer', 795634069.0, 2004, 142, 26, 21, 'фильм');


INSERT INTO filmgenres VALUES 
    (1, 2),
    (1, 3),
    (2, 2),
    (3, 2),
    (4, 2),
    (5, 2),
    (5, 9),
    (5, 3),
    (6, 9),
    (6, 10),
    (7, 9),
    (7, 10),
    (8, 9),
    (8, 10);


INSERT INTO countryproduction VALUES 
    (1, 1),
    (1, 2),
    (2, 1),
    (3, 4),
    (4, 4),
    (5, 4),
    (5, 5),
    (6, 3),
    (6, 4),
    (7, 3),
    (7, 4),
    (8, 3),
    (8, 4);


INSERT INTO filmcast VALUES 
    (1, 3), (1, 4), (1, 5),
    (2, 8), (2, 9), (2, 10),
    (3, 8), (3, 13),
    (4, 8), (4, 16),
    (5, 8), (5, 19), (5, 20),
    (6, 23), (6, 24), (6, 25),
    (7, 23), (7, 24), (7, 25),
    (8, 23), (8, 24), (8, 25);


INSERT INTO review VALUES 
    (DEFAULT, 1, 1, 'норм', 2, 5.0, current_date),
    (DEFAULT, 2, 1, ')', 2, 6.0, current_date),
    (DEFAULT, 3, 1, 'ао ', 3, 7.0, current_date),
    (DEFAULT, 4, 1, 'неплохо', 3, 7.0, current_date),
    (DEFAULT, 5, 1, 'скучно', 1, 2.0, current_date),
    (DEFAULT, 6, 1, 'да', 3, 7.0, current_date),
    (DEFAULT, 7, 1, 'интересно', 3, 8.0, current_date),
    (DEFAULT, 8, 1, 'отвал башки', 3, 10.0, current_date),
    (DEFAULT, 1, 2, 'ввввввв', 1, 1.0, current_date),
    (DEFAULT, 8, 2, '', 0, 9.0, current_date),
    (DEFAULT, 8, 3, '', 0, 5.0, current_date),
    (DEFAULT, 8, 4, '', 0, 10.0, current_date),
    (DEFAULT, 8, 5, 'ffff', 1, 0, current_date);

INSERT INTO recommended VALUES
    (1, 2),
    (1, 3),
    (2, 1),
    (3, 1),
    (4, 3),
    (3, 4),
    (5, 6),
    (6, 5),
    (6, 7),
    (6, 8),
    (7, 6),
    (7, 8),
    (8, 6),
    (8, 7);

INSERT INTO bookmark VALUES
     (1, 1),
     (2, 1),
     (3, 1),
     (5, 1),
     (4, 2),
     (3, 2),
     (1, 3),
     (2, 3),
     (7, 3),
     (8, 4),
     (2, 5),
     (1, 6),
     (4, 7),
     (3, 8);

INSERT INTO subscription VALUES
     (1, 1),
     (2, 1),
     (3, 1),
     (5, 1),
     (4, 2),
     (3, 2),
     (1, 3),
     (2, 3),
     (7, 3),
     (8, 4),
     (2, 5),
     (1, 6),
     (4, 7),
     (3, 8);
           /*
SELECT film.*, p.person_id,p.name_en,p.name_rus,p.picture_url,p.career, p1.person_id,p1.name_en,p1.name_rus,p1.picture_url,p1.career FROM film JOIN person p on film.director = p.person_id JOIN person p1 on film.screenwriter = p1.person_id WHERE film_id = 7
SELECT country.country_name FROM country JOIN countryproduction c ON country.country_id = c.country_id WHERE c.Film_ID = 7;
SELECT genre.* FROM genre JOIN filmgenres f on genre.genre_id = f.genre_id WHERE f.film_id = 7;
SELECT p.person_id,p.name_en,p.name_rus,p.picture_url,p.career  FROM person p JOIN filmcast f on p.person_id = f.person_id WHERE f.film_id = 7;
SELECT * FROM review
SELECT COUNT(*) FROM Review WHERE Film_ID = 6   ;
SELECT review.*, p.first_name, p.surname FROM review join profile p on p.user_id = review.author_id WHERE film_id = 6 LIMIT 10 OFFSET 0
SELECT f.film_id, f.title, f.poster_url FROM recommended r join film f on f.film_id = r.recommended_id WHERE r.film_id = 7
SELECT * FROM review where film_id = 2                           ;
SELECT AVG(stars) FROM review WHERE film_id =4 AND (NOT type = 0);

SELECT review.* FROM review JOIN  WHERE film_id = 1 AND author_id = 1;
            */

-- SELECT review.*, p.first_name, p.surname FROM review join profile p on p.user_id = review.author_id WHERE film_id = 8 AND (NOT type = 0) LIMIT 10 OFFSET 0
-- SELECT review.*, p.first_name, p.surname, p.picture_url FROM review join profile p on p.user_id = review.author_id WHERE film_id = 8 AND (NOT type = 0) LIMIT 10;
-- select * from review
