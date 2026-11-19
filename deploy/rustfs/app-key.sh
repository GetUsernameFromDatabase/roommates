#!/bin/sh
# Gives the app its own RustFS user, limited to the bucket S3_BUCKET of .env, see deploy/README.md 7.6
#
#   sh deploy/rustfs/app-key.sh          create the user, or re-apply the policy below after an edit
#   sh deploy/rustfs/app-key.sh rotate   give the user a new secret, the old one stops working
#
# Prints APP_S3_ACCESS_KEY / APP_S3_SECRET_KEY for .env whenever it makes a secret.
# Needs this checkout's rustfs running: docker compose up -d --wait rustfs
# Admin API: https://docs.rustfs.com/administration/iam/policies
set -eu

# on the host: run this same file inside the rustfs container of the checkout it lives in,
# which has the root keys and S3_BUCKET in its env, so nothing reads .env
if [ "${1:-}" != --in-container ]; then
  case ${1:-} in
    '' | rotate) ;;
    *) echo "usage: sh deploy/rustfs/app-key.sh [rotate]" >&2; exit 2 ;;
  esac
  cd "$(dirname "$0")/../.."
  exec docker compose exec -T rustfs sh -s -- --in-container "$@" < deploy/rustfs/app-key.sh
fi
rotate=${2:-}

user=roommates-app
bucket=${S3_BUCKET:?S3_BUCKET is not set in .env}

# what the app does, see app/internal/storage/storage.go (CreateBucket: only when it is missing)
policy=$(cat <<JSON
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["s3:ListBucket", "s3:CreateBucket"],
      "Resource": ["arn:aws:s3:::$bucket"]
    },
    {
      "Effect": "Allow",
      "Action": ["s3:GetObject", "s3:PutObject", "s3:DeleteObject"],
      "Resource": ["arn:aws:s3:::$bucket/*"]
    }
  ]
}
JSON
)

# curl signs like an S3 client: provider aws:amz, region (RustFS ignores it), service s3;
# curl needs all of them, it cannot read region and service from "localhost"
sigv4=aws:amz:us-east-1:s3

# admin METHOD PATH [BODY]: prints the HTTP status, the response body goes to /tmp/app-key.out
admin() {
  curl -sS -o /tmp/app-key.out -w '%{http_code}' -X "$1" --aws-sigv4 "$sigv4" \
    --user "$RUSTFS_ACCESS_KEY:$RUSTFS_SECRET_KEY" --data-binary "${3:-}" \
    "http://localhost:9000/rustfs/admin/v3/$2"
}

# must METHOD PATH [BODY]: admin, stopping unless it returns 200
must() {
  code=$(admin "$@")
  [ "$code" = 200 ] || { echo "$1 $2 failed ($code): $(cat /tmp/app-key.out)" >&2; exit 1; }
}

new_secret() {
  secret=$(head -c 32 /dev/urandom | od -An -tx1 | tr -d ' \n' | tr a-f A-F)
  must PUT "add-user?accessKey=$user" "{\"secretKey\":\"$secret\",\"status\":\"enabled\"}"
}

must PUT "add-canned-policy?name=$user" "$policy"
secret=
code=$(admin GET "user-info?accessKey=$user")
case $code in
  404) new_secret ;;
  200) if [ "$rotate" = rotate ]; then new_secret; fi ;;
  *) echo "user-info failed ($code): $(cat /tmp/app-key.out)" >&2; exit 1 ;;
esac
must PUT "set-user-or-group-policy?policyName=$user&userOrGroup=$user&isGroup=false"
rm -f /tmp/app-key.out

if [ -n "$secret" ]; then
  echo "# bucket $bucket, paste into .env, then: docker compose up -d"
  echo "APP_S3_ACCESS_KEY=$user"
  echo "APP_S3_SECRET_KEY=$secret"
else
  echo "Policy for bucket $bucket applied to $user, secret unchanged (new one: sh deploy/rustfs/app-key.sh rotate)."
fi
