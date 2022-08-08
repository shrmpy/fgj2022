# fgj2022
For myFirst Game Jam entry

Summer 2022 Heal theme of
 [My First Game Jam](https://itch.io/jam/my-first-game-jam-summer-2022)

When done 🤞, download will be on the
 [itch.io site here](https://shrmpy.itch.io/fgj2022)


## Quickstart
```bash
git clone https://github.com/shrmpy/fgj2022
cd fgj2022 && go build 
./fgj2022
```
## Build in Local Container
```bash
cd fgj2022
docker build -t bc .
docker run -ti --rm --entrypoint sh -v $PWD:/opt/test bc
go build -o test
cp test /opt/test/fgj2022
exit
./fgj2022
```
## Make your own snap package
[![fgj2022](https://snapcraft.io/fgj2022/badge.svg)](https://snapcraft.io/fgj2022)
```bash
# ub server includes a empty lxd?
sudo snap remove --purge lxd
# reinstall lxd
sudo snap install lxd
sudo lxd init --auto
sudo usermod -a -G lxd ${USER}
# view config
lxc version
lxc profile show default
lxc storage show default
echo 'export SNAPCRAFT_BUILD_ENVIRONMENT=lxd' >> ~/.profile
sudo reboot
# retrieve YAML 
git clone https://github.com/shrmpy/fgj2022.git
cd fgj2022
# make snap 
snapcraft
# local install
sudo snap install fgj2022_0.0.6_arm64.snap --dangerous
# start program
fgj2022
```


## Credits

Wazero imports
 by [Yeicor](https://github.com/Yeicor/sdfx-isosurface)

Wazero embed
 by [Fernando Talavera](https://github.com/efejjota/ebiten-wasm-graphics)

Wazero
 by [Tetrate.io](https://github.com/tetratelabs/wazero) ([LICENSE](https://github.com/tetratelabs/wazero/blob/main/LICENSE))

WASI flite 
 by [Jakub Konka](http://www.jakubkonka.com/2019/04/20/wasi-flite.html)

Flite bindings
 by [Milan Nikolic](https://github.com/gen2brain/flite-go) ([LICENSE](https://github.com/gen2brain/flite-go/blob/master/LICENSE))

Github workflow
 by [Siôn le Roux](https://github.com/sinisterstuf/ebiten-game-template) ([LICENSE](https://github.com/sinisterstuf/ebiten-game-template/blob/main/LICENSE))

Font Renderer
 by [tinne26](https://github.com/tinne26/etxt)
 ([LICENSE](https://github.com/tinne26/etxt/blob/main/LICENSE))

Ebitengine
 by [Hajime Hoshi](https://github.com/hajimehoshi/ebiten/)
 ([LICENSE](https://github.com/hajimehoshi/ebiten/blob/main/LICENSE))

DejaVu Sans Mono
 by [DejaVu](https://dejavu-fonts.github.io/)
 ([LICENSE](https://github.com/dejavu-fonts/dejavu-fonts/blob/master/LICENSE))

