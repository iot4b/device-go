#!/bin/sh

# Пути к репозиториям и сборкам
REPO_PATH="/opt/iot4b/repo/packages"
BUILD_PATH="./builds"

# Ассоциативный массив: Архитектура -> Имя пакета
declare -A ARCH_FILES=(
    ["mipsel"]="iot4b_openwrt.ipk"
    ["mipsel-3.4_kn"]="iot4b_keenetic.ipk"
    #["armv7l"]="iot4b_openwrtarmv7l.ipk"
    #["aarch64"]="iot4b_openwrtaarch64.ipk"
)

# Очистка старых пакетов
for ARCH in "${!ARCH_FILES[@]}"; do
    rm -f "$REPO_PATH/$ARCH/"*
done

# Копирование новых пакетов
for ARCH in "${!ARCH_FILES[@]}"; do
    cp "$BUILD_PATH/${ARCH_FILES[$ARCH]}" "$REPO_PATH/$ARCH/iot4b-$ARCH.ipk"
done

# Создание индексов
for ARCH in "${!ARCH_FILES[@]}"; do
    cd "$REPO_PATH/$ARCH" || exit 1
    opkg-make-index -a ./ > Packages
    gzip -k Packages
done