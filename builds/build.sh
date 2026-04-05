#!/bin/bash

# Прекращает выполнение скрипта при ошибке
set -e

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
eval "$("${ROOT_DIR}/scripts/release-vars.sh")"

# Определение переменных
PACKAGE_NAME="iot4b"
VERSION="${IOT4B_VERSION}"
ARCH="mips_siflower"
MAINTAINER="IOT4B <sp@golbex.com>"
DESCRIPTION="IOT4B device client"
OS="$1"
ARCH="$2"
LD_FLAGS="-s -w -X device-go/packages/buildinfo.Version=${IOT4B_VERSION} -X device-go/packages/buildinfo.Commit=${IOT4B_COMMIT} -X device-go/packages/buildinfo.BuildDate=${IOT4B_BUILD_DATE}"

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
cd "${ROOT_DIR}" || exit 1
env GOOS=linux "${GOARCH_FLAGS[@]}" go build -ldflags="${LD_FLAGS}" -o "${BUILD_DIR}/opt/${PACKAGE_NAME}/${PACKAGE_NAME}" ./main.go
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
