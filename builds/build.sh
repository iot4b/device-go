#!/bin/bash

# Прекращает выполнение скрипта при ошибке
set -e

# Определение переменных
PACKAGE_NAME="iot4b"
VERSION="1.2.6"
ARCH="mips_siflower"
MAINTAINER="IOT4B <sp@golbex.com>"
DESCRIPTION="IOT4B device client"
OS="$1"
ARCH="$2"

if [[ "$ARCH" == "armv7l" ]]; then
  GOARCH_FLAGS=(GOARCH=arm GOARM=7)
elif [[ "$ARCH" == "mips"* ]]; then
  GOARCH_FLAGS=(GOARCH=mipsle)
else
  GOARCH_FLAGS=(GOARCH=arm64)
fi

if [[ $OS == "keenetic" ]]; then
  DEPENDS="ndmq, libndm"
else
  DEPENDS="libc"
fi

BUILD_DIR="./builds/${PACKAGE_NAME}_${OS}_${ARCH}"
IPK_PATH="../${PACKAGE_NAME}_${OS}_${ARCH}.ipk"

# Создание файла control с правильным Installed-Size
INSTALLED_SIZE=$(du -sk "${BUILD_DIR}/opt" | awk '{print $1}')

cat <<EOF > "${BUILD_DIR}/control"
Package: ${PACKAGE_NAME}
Version: ${VERSION}
Depends: ${DEPENDS}
Source: 
Section: utils
License: GPL-2.0+
Architecture: ${ARCH}
Installed-Size: ${INSTALLED_SIZE}
Maintainer: ${MAINTAINER}
Description: ${DESCRIPTION}
EOF

# Компиляция бинарника Go для mipsle
echo "Компиляция бинарника Go для ${OS} ${ARCH}..."
env GOOS=linux "${GOARCH_FLAGS[@]}" go build -ldflags="-s -w" -o "${BUILD_DIR}/opt/${PACKAGE_NAME}/${PACKAGE_NAME}" main.go
echo "Бинарник собран. Размер:"
ls -lh "${BUILD_DIR}/opt/${PACKAGE_NAME}/${PACKAGE_NAME}" | awk '{print "Размер бинарника:", $5}'

# Переход в директорию сборки
cd "${BUILD_DIR}" || exit 1

# Удаление старого ipk, если существует
rm -f "${IPK_PATH}"

# Создание control.tar.gz
echo "Создание control.tar.gz..."
tar --format=ustar -czvf control.tar.gz control postinst postrm

# Создание data.tar.gz
echo "Создание data.tar.gz..."
if [[ $OS == "keenetic" ]]; then
  tar --format=ustar -czvf data.tar.gz opt
else
  tar --format=ustar -czvf data.tar.gz opt etc
fi

# Создание debian-binary
echo "Создание debian-binary..."
echo "2.0" > debian-binary

# Сборка .ipk пакета
echo "Сборка .ipk пакета..."

if [[ $OS == "keenetic" ]]; then
  tar --format=ustar -czvf "${IPK_PATH}" debian-binary control.tar.gz data.tar.gz
else
  ar rcs "${IPK_PATH}" debian-binary control.tar.gz data.tar.gz
fi

echo "Пакет создан: ${IPK_PATH}"
ls -lh "${IPK_PATH}" | awk '{print "Размер IPK пакета:", $5}'

# Очистка временных файлов
echo "Очистка временных файлов..."
rm -f control.tar.gz data.tar.gz debian-binary
rm -f opt/${PACKAGE_NAME}/${PACKAGE_NAME}

echo "Сборка и упаковка завершены успешно!"
echo ""