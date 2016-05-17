## go_rest_api

Trying to figure out authentication and interactions with a database in go.

+-------------+
| client      |
+-------------+
    |
    | username/password
    |
+--------+                           +--------+
|        |      username/password    |        |
| client |   --------------------->  | server |
|        |                           |        |
|        |    JSON Web Token         |        |
|        |  <---------------------   |        |
+--------+                           +--------+
