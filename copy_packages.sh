#!/bin/sh

# Пути к репозиторию и сборкам
REPO_PATH="/opt/iot4b/repo/packages"
BUILD_PATH="./builds"

# Функция для обработки пакетов
update_repo() {
    ARCH=$1
    PKG_FILE=$2

    # Проверяем и создаём папку
    if [ ! -d "$REPO_PATH/$ARCH" ]; then
        echo "Создаю папку: $REPO_PATH/$ARCH"
        mkdir -p "$REPO_PATH/$ARCH"
    fi

    # Удаляем старые файлы
    echo "Удаляю старые файлы в $REPO_PATH/$ARCH"
    rm -f "$REPO_PATH/$ARCH/"*

    # Копируем новый пакет
    SRC_FILE="$BUILD_PATH/$PKG_FILE"
    DEST_FILE="$REPO_PATH/$ARCH/iot4b-$ARCH.ipk"

    if [ -f "$SRC_FILE" ]; then
        echo "Копирую $SRC_FILE -> $DEST_FILE"
        cp "$SRC_FILE" "$DEST_FILE"
    else
        echo "Ошибка: Файл $SRC_FILE не найден!"
        exit 1
    fi

    # Создание индексов
    echo "Создаю индекс в $REPO_PATH/$ARCH"
    cd "$REPO_PATH/$ARCH" || exit 1
    opkg-make-index -a ./ > Packages
    gzip -k Packages
}

# Вызов функции 4 раза для разных архитектур
update_repo "mipsel" "iot4b_openwrt.ipk"
update_repo "mipsel-3.4_kn" "iot4b_keenetic.ipk"
update_repo "armv7l" "iot4b_keenetic.ipk"
update_repo "aarch64" "iot4b_keenetic.ipk"

