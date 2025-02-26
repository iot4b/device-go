#!/bin/sh

# Пути к репозиторию и сборкам
REPO_PATH="/opt/iot4b/repo/packages"
BUILD_PATH="./builds"

# Функция для обработки пакетов
update_repo() {
    ARCH=$1
    PKG_FILE=$2
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
    DEST_FILE="$ARCH_PATH/iot4b-$ARCH.ipk"

    if [ -f "$SRC_FILE" ]; then
        echo "Копирую $SRC_FILE -> $DEST_FILE"
        cp "$SRC_FILE" "$DEST_FILE"
    else
        echo "Ошибка: Файл $SRC_FILE не найден!"
        exit 1
    fi

    # Создание индексов
    echo "Создаю индекс в $ARCH_PATH"
    opkg-make-index -a "$ARCH_PATH" > "$ARCH_PATH/Packages"
    gzip -k "$ARCH_PATH/Packages"
}

# Вызов функции 4 раза для разных архитектур
update_repo "armv7l" "iot4b_openwrtarmv7l.ipk"
update_repo "mipsel" "iot4b_openwrt.ipk"
update_repo "aarch64" "iot4b_openwrtaarch64.ipk"
update_repo "mipsel-3.4_kn" "iot4b_keenetic.ipk"


