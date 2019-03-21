#!/bin/dash

./clean.sh
name=$(cat ./PKGBUILD | grep pkgname= | cut -b 9-)
version=$(cat ./PKGBUILD | grep pkgver= | cut -b 8-)
echo "$name-$version release"

makepkg --printsrcinfo > .SRCINFO
rm -rf /tmp/iuqyeiuqdkjsahdjakshd
mkdir /tmp/iuqyeiuqdkjsahdjakshd
cp ./PKGBUILD ./.SRCINFO /tmp/iuqyeiuqdkjsahdjakshd
cd /tmp/iuqyeiuqdkjsahdjakshd
git init
git add -vA
git remote add aur-nelc-git ssh://aur@aur.archlinux.org/nelc-git.git
#git log -p master aur-nelc-git/master
git commit -m "$name-$version release"
git fetch aur-nelc-git

