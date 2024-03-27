NODE_OPTIONS=--openssl-legacy-provider npm run build:h5
(cd ..; zip frontend_build -r frontend/dist/build/h5 -9; rm -rf frontend/dist)
# scp ../frontend_build.zip art:~/fhl-web
# server: (cd ~/fhl-web; rm -rf frontend/dist; unzip frontend_build.zip; rm frontend_build.zip)
