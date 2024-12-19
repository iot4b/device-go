include $(TOPDIR)/rules.mk

PKG_VERSION:=1.0.0
PKG_RELEASE:=1

PKG_SOURCE:=$(PKG_VERSION).tar.gz
PKG_SOURCE_URL:=https://github.com/ever-iot/device-go/archive/refs/tags
PKG_HASH:=skip

PKG_BUILD_DIR:=$(BUILD_DIR)$(PKG_VERSION)

include $(INCLUDE_DIR)/package.mk

define Package/iot4b
  SECTION:=utils
  CATEGORY:=Utilities
  TITLE:=iot4b
  DEPENDS:=
  MAINTAINER:=ServicePusher <sp@golbex.com>
endef

define Package/iot4b/description
IOT4B module for iot device.
endef

define Package/iot4b/install
	$(INSTALL_DIR) $(1)/usr/bin
	$(INSTALL_BIN) $(PKG_BUILD_DIR)/iot4b $(1)/usr/bin/iot4b

	$(INSTALL_DIR) $(1)/etc/config
	$(INSTALL_CONF) ./files/iot4b.config $(1)/etc/config/iot4b

    $(INSTALL_DIR) $(1)/etc/config
	$(INSTALL_CONF) ./files/iot4b.init.json $(1)/etc/config/iot4b/iot4b.init.json

    $(INSTALL_DIR) $(1)/etc/config
	$(INSTALL_CONF) ./files/iot4b.yml $(1)/etc/config/iot4b/iot4b.yml

	$(INSTALL_DIR) $(1)/etc/init.d
	$(INSTALL_BIN) ./files/iot4b.init $(1)/etc/init.d/iot4b
endef

$(eval $(call BuildPackage,iot4b))