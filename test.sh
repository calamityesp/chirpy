#! /bin/bash


echo "creating user key@gmail.com, with password 123456"
curl -X POST \
-H "content-type: application/json" \
-d '{"email":"key@gmail.com", "password":"123456"}' localhost:8080/api/users

sleep 2
echo "logging in and generating hashed password"
curl -X POST \
   -H "content-type: application/json" \
   -d '{"email":"key@gmail.com", "password":"123456"}' \
   localhost:8080/api/login

sleep 2
echo -n "Test webHook? "
read webhook

if [[ "${webhook,,}" == "y"* ]]; then
   curl -X POST \
      -H "content-type: application/json" \
      -d '{"data": {"user_id":1}, "event":"user.upgraded"}' \
      localhost:8080/api/polka/webhooks
fi

