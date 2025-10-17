#!/bin/sh

# Пути к репозиторию и сборкам
REPO_PATH="/opt/iot4b/repo/packages"
BUILD_PATH="./builds"

# Функция для обработки пакетов
update_repo() {
    OS=$1
    ARCH=$2
    PKG_FILE="iot4b_${OS}_${ARCH}.ipk"
    ARCH_PATH="$REPO_PATH/$ARCH"

    # Проверяем и создаём папку
    if [ ! -d "$ARCH_PATH" ]; then
        echo "Создаю папку: $ARCH_PATH"
        mkdir -p "$ARCH_PATH"
    fi

    # Удаляем старые файлы
    echo "Удаляю старые файлы в $ARCH_PATH"
    rm -f "$ARCH_PATH/"*

    # Копируем новый пакет
    SRC_FILE="$BUILD_PATH/$PKG_FILE"
    DEST_FILE="$ARCH_PATH/iot4b-$OS-$ARCH.ipk"

    if [ -f "$SRC_FILE" ]; then
        echo "Копирую $SRC_FILE -> $DEST_FILE"
        cp "$SRC_FILE" "$DEST_FILE"
    else
        echo "Ошибка: Файл $SRC_FILE не найден!"
        exit 1
    fi

    # Создание индексов
    echo "Создаю индекс в $ARCH_PATH"
    if file "$DEST_FILE" | grep -q "gzip compressed data"; then
        opkg_make_index "$OS" "$ARCH"
    else
        opkg-make-index -a "$ARCH_PATH" > "$ARCH_PATH/Packages"
    fi
    gzip -k "$ARCH_PATH/Packages"
}

# alternative to opkg-make-index utility for tar based ipk packages
opkg_make_index() {
    OS=$1
    ARCH=$2
    ARCH_PATH="$REPO_PATH/$ARCH"
    PKG_FILE="$ARCH_PATH/iot4b-$OS-$ARCH.ipk"
    CONTROL_DIR=$(mktemp -d)

    # Extract control info
    tar -xzf "$PKG_FILE" -C "$CONTROL_DIR"
    tar -xzf "$CONTROL_DIR/control.tar.gz" -C "$CONTROL_DIR"

    # Collect info
    MD5SUM=$(md5sum "$PKG_FILE" | cut -d ' ' -f1)
    SIZE=$(stat --printf="%s" "$PKG_FILE")
    FILENAME=$(basename "$PKG_FILE")

    # Read fields from control file
    PACKAGE=$(grep "^Package:" "$CONTROL_DIR/control" | cut -d ' ' -f2)
    VERSION=$(grep "^Version:" "$CONTROL_DIR/control" | cut -d ' ' -f2)
    SECTION=$(grep "^Section:" "$CONTROL_DIR/control" | cut -d ' ' -f2)
    ARCHITECTURE=$(grep "^Architecture:" "$CONTROL_DIR/control" | cut -d ' ' -f2)
    MAINTAINER=$(grep "^Maintainer:" "$CONTROL_DIR/control" | cut -d ' ' -f2-)
    DESCRIPTION=$(grep "^Description:" "$CONTROL_DIR/control" | cut -d ' ' -f2-)

    # Generate Packages file
    cat << EOF > "$ARCH_PATH/Packages"
Package: $PACKAGE
Version: $VERSION
Section: $SECTION
Architecture: $ARCHITECTURE
Maintainer: $MAINTAINER
MD5Sum: $MD5SUM
Size: $SIZE
Filename: $FILENAME
Description: $DESCRIPTION
EOF

    rm -rf "$CONTROL_DIR"
}

update_apt() {
  APT_REPO_PATH="/opt/iot4b/repo/apt"
  APT_BUILD_PATH="$(cd "$(dirname "$0")" >/dev/null 2>&1 && pwd -P)/apt"

  echo "Обновляю APT репозиторий"
  cp "${APT_BUILD_PATH}/iot4bd_amd64.deb" "${APT_REPO_PATH}/iot4bd_amd64.deb"
  cd ${APT_REPO_PATH} || exit 1
  dpkg-scanpackages --arch amd64 . > dists/stable/main/binary-amd64/Packages
  cat dists/stable/main/binary-amd64/Packages | gzip -9 > dists/stable/main/binary-amd64/Packages.gz
  cd dists/stable || exit 1
  "${APT_BUILD_PATH}/generate-release.sh" > Release
  cat Release | gpg --default-key iot4bd -abs > Release.gpg
  cat Release | gpg --default-key iot4bd -abs --clearsign > InRelease
}

# Вызов функции 4 раза для разных архитектур
update_repo "keenetic" "mipsel-3.4_kn"
update_repo "keenetic" "aarch64-3.10_kn"
update_repo "openwrt" "mips_siflower"
update_repo "openwrt" "armv7l"
update_repo "openwrt" "aarch64"

# копируем файлы в папку REPO_PATH
cp "${BUILD_PATH}/iot4b_install.sh" "${REPO_PATH}/install.sh"

cp "${BUILD_PATH}/libndm_1.8.0-1_mipsel-3.4_kn.ipk" "${REPO_PATH}/mipsel-3.4_kn/libndm_1.8.0-1_mipsel-3.4_kn.ipk"
cp "${BUILD_PATH}/ndmq_1.0.2-7_mipsel-3.4_kn.ipk" "${REPO_PATH}/mipsel-3.4_kn/ndmq_1.0.2-7_mipsel-3.4_kn.ipk"

cp "${BUILD_PATH}/libndm_1.1.25-1_aarch64-3.10_kn.ipk" "${REPO_PATH}/aarch64-3.10_kn/libndm_1.1.25-1_aarch64-3.10_kn.ipk"
cp "${BUILD_PATH}/ndmq_1.0.2-11_aarch64-3.10_kn.ipk" "${REPO_PATH}/aarch64-3.10_kn/ndmq_1.0.2-11_aarch64-3.10_kn.ipk"

update_apt
