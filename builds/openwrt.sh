#!/bin/bash

# Прекращает выполнение скрипта при ошибке
set -e

# Определение переменных
PACKAGE_NAME="iot4b"
VERSION="1.0"
ARCH="mips_siflower"
MAINTAINER="IOT4B <sp@golbex.com>"
DESCRIPTION="IOT4B device client"

if [[ "$1" == "arm" ]]; then
  ARCH="armv7l"
  GOARCH_FLAGS=(GOARCH=arm GOARM=7)
elif [[ "$1" == "arm64" ]]; then
  ARCH="aarch64"
  GOARCH_FLAGS=(GOARCH=arm64)
else
  ARCH="mips_siflower"
  GOARCH_FLAGS=(GOARCH=mipsle)
fi

BUILD_DIR="./builds/${PACKAGE_NAME}_openwrt_${ARCH}"
IPK_PATH="../${PACKAGE_NAME}_openwrt_${ARCH}.ipk"

# Создание файла control с правильным Installed-Size
INSTALLED_SIZE=$(du -sk "${BUILD_DIR}/opt" | awk '{print $1}')

cat <<EOF > "${BUILD_DIR}/control"
Package: ${PACKAGE_NAME}
Version: ${VERSION}
Depends: libc
Source: 
Section: utils
License: GPL-2.0+
Architecture: ${ARCH}
Installed-Size: ${INSTALLED_SIZE}
Maintainer: ${MAINTAINER}
Description: ${DESCRIPTION}
EOF

# Компиляция бинарника Go для mipsle
echo "Компиляция бинарника Go для ${ARCH}..."
env GOOS=linux "${GOARCH_FLAGS[@]}" go build -ldflags="-s -w" -o "${BUILD_DIR}/opt/${PACKAGE_NAME}/${PACKAGE_NAME}" main.go
echo "Бинарник собран. Размер:"
ls -lh "${BUILD_DIR}/opt/${PACKAGE_NAME}/${PACKAGE_NAME}" | awk '{print "Размер бинарника:", $5}'

# Переход в директорию сборки
cd "${BUILD_DIR}" || exit 1

# Удаление старого ipk, если существует
rm -f "${IPK_PATH}"

# Создание control.tar.gz
echo "Создание control.tar.gz..."
tar --format=ustar -czf control.tar.gz control postinst postrm

# Создание data.tar.gz
echo "Создание data.tar.gz..."
tar --format=ustar -czf data.tar.gz opt etc

# Создание debian-binary
echo "Создание debian-binary..."
echo "2.0" > debian-binary

# Сборка .ipk пакета
echo "Сборка .ipk пакета..."

#tar --format=ustar -czvf "${IPK_PATH}" debian-binary control.tar.gz data.tar.gz
ar rcs  "${IPK_PATH}" debian-binary control.tar.gz data.tar.gz

echo "Пакет создан: ${IPK_PATH}"
ls -lh "${IPK_PATH}" | awk '{print "Размер IPK пакета:", $5}'

# Очистка временных файлов
echo "Очистка временных файлов..."
rm -f control.tar.gz data.tar.gz debian-binary
rm -f opt/${PACKAGE_NAME}/${PACKAGE_NAME}

echo "Сборка и упаковка завершены успешно!"