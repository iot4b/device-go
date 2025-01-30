#!/bin/bash

# Прекращает выполнение скрипта при ошибке
set -e

# Определение переменных
PACKAGE_NAME="iot4b"
VERSION="1.0"
ARCH="mips_siflower"
MAINTAINER="IOT4B <sp@golbex.com>"
DESCRIPTION="IOT4B device client"
BUILD_DIR="./builds/${PACKAGE_NAME}_openwrt2"
IPK_PATH="../${PACKAGE_NAME}_openwrt2.ipk"


# Создание файла control с правильным Installed-Size
INSTALLED_SIZE=$(du -sk "${BUILD_DIR}/opt" | awk '{print $1}')

cat <<EOF > "${BUILD_DIR}/control"
Package: ${PACKAGE_NAME}
Version: ${VERSION}
Depends: libc
Source: package/libs/iot4b
Section: utils
License: GPL-2.0+
Architecture: ${ARCH}
Installed-Size: ${INSTALLED_SIZE}
Maintainer: ${MAINTAINER}
Description: ${DESCRIPTION}
EOF

# Компиляция бинарника Go для mipsle
echo "Компиляция бинарника Go для ${ARCH}..."
GOOS=linux GOARCH=mipsle go build -ldflags="-s -w" -o "${BUILD_DIR}/opt/${PACKAGE_NAME}/${PACKAGE_NAME}" main.go
echo "Бинарник собран. Размер:"
ls -lh "${BUILD_DIR}/opt/${PACKAGE_NAME}/${PACKAGE_NAME}" | awk '{print "Размер бинарника:", $5}'

# Переход в директорию сборки
cd "${BUILD_DIR}" || exit 1

# Удаление старого ipk, если существует
rm -f "${IPK_PATH}"

# Создание control.tar.gz
echo "Создание control.tar.gz..."
tar -czf control.tar.gz control postinst

# Создание data.tar.gz
echo "Создание data.tar.gz..."
tar -czf data.tar.gz opt etc

# Создание debian-binary
echo "Создание debian-binary..."
echo "2.0" > debian-binary

# Сборка .ipk пакета
echo "Сборка .ipk пакета..."
tar -czvf "${IPK_PATH}" debian-binary control.tar.gz data.tar.gz

echo "Пакет создан: ${IPK_PATH}"
ls -lh "${IPK_PATH}" | awk '{print "Размер IPK пакета:", $5}'

# Очистка временных файлов
echo "Очистка временных файлов..."
rm -f control.tar.gz data.tar.gz debian-binary control
rm -f opt/${PACKAGE_NAME}/${PACKAGE_NAME}

echo "Сборка и упаковка завершены успешно!"