# Buter - Brute Force Tool

Works as other tools but multiple payload set can \
be used.

Example 1: 
```buter.exe -b '{\"email\":\"test@gmail.com\",\"password\":\"!abc!\"}' -h '{\"Cookie\": \"!some_value!\"}' -m POST -u "https://bla.bla.com?param=!some_value!" -p /path/payload-list1 -p /path/payload-list2 -p /path/payload-list3 -a cluster```

Example 2: 
```buter.exe -b '{\"email\":\"test@gmail.com\",\"password\":\"!abc!\"}' -m POST -u "https://bla.bla.com/auth/login" -p /path/payload-list1 -a sniper```