#!/bin/bash
echo "Recreating rpmbuild directory"
rm -rvf /root/rpmbuild/
rpmdev-setuptree
echo "Building SRPM"
rpmbuild --undefine=_disable_source_fetch -bs /project/gpt-telegram-bot/rpm/gpt-telegram-bot.spec
mkdir -p ~/.config
mv /project/gpt-telegram-bot/copr ~/.config/copr
copr-cli build gpt-telegram-bot /root/rpmbuild/SRPMS/*.src.rpm
