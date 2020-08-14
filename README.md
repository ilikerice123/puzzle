# Puzzle

Multiplayer puzzle server & client similar to [epuzzle.info](http://epuzzle.info)

HTTP API:
- `/api/images/...`
  - This is the route used to serve static images. Each image uploaded gets put in a unique folder, 
so that it can accessed like: 
    - `/api/images/<uuid>/original.jpeg`
    - `/api/images/<uuid>/preview.jpeg` (scaled down version to 200px length)
    - `/api/images/<uuid>/original_Y_X.jpeg` for pieces

WSS API:


TODO:
- have logging
- create reaper that goes through and cleans up directory/user data