#!/bin/sh

# This script is for installing the latest version of Turso CLI on your machine.

set -e

# Terminal ANSI escape codes.
reset="\033[0m"
bright_blue="${reset}\033[34;1m"

probe_arch() {
    ARCH=$(uname -m)
    case $ARCH in
        x86_64) ARCH="x86_64"  ;;
        aarch64) ARCH="arm64" ;;
        arm64) ARCH="arm64" ;;
        *) printf "Architecture ${ARCH} is not supported by this installation script\n"; exit 1 ;;
    esac
}

probe_os() {
    OS=$(uname -s)
    case $OS in
        Darwin) OS="Darwin" ;;
        Linux) OS="Linux" ;;
        *) printf "Operating system ${OS} is not supported by this installation script\n"; exit 1 ;;
    esac
}

print_logo() {
    printf "${bright_blue}
                 .:                                 .:
  .\$\$.   \$\$:   .\$\$\$:                                \$\$\$^    \$\$:   ~\$^
  .\$\$\$!:\$\$\$  .\$\$\$\$~                                 .\$\$\$\$^  !\$\$~^\$\$\$~
    \$\$\$\$\$\$ .\$\$\$\$\$~                                   .\$\$\$\$\$^ \$\$\$\$\$\$:
     !\$\$\$\$\$\$\$\$\$\$~                                     .\$\$\$\$\$\$\$\$\$\$\$
      :\$\$\$\$\$\$\$\$~                                       .\$\$\$\$\$\$\$\$!
     .\$\$\$\$\$\$\$\$~                                         .\$\$\$\$\$\$\$\$^
    .\$\$\$\$\$\$\$\$!       ~\$!                       :\$\$.      :\$\$\$\$\$\$\$\$^
     \$\$\$\$\$\$\$\$\$\$\$!^::\$\$\$\$\$^...................:\$\$\$\$\$!.^~\$\$\$\$\$\$\$\$\$\$\$:
     \$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$
     :\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$!
        :^!\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$~:.
           :\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$!
      \$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$:
      :\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$~
        ^\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$~.
           :\$\$\$\$\$:   .^~!\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$!^:.   \$\$\$\$\$!
           :\$\$\$\$\$!.         .!\$\$\$\$\$\$\$\$\$\$\$\$.         .^\$\$\$\$\$!
           :\$\$\$\$\$\$\$\$\$\$!^:.   ~\$\$\$\$\$\$\$\$\$\$\$\$    .^~\$\$\$\$\$\$\$\$\$\$!
           :\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$. ~\$\$\$\$\$\$\$\$\$\$\$\$  \$\$\$\$\$\$\$\$\$\$\$\$\$\$\$!
           :\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$: ~\$\$\$\$\$\$\$\$\$\$\$\$  \$\$\$\$\$\$\$\$\$\$\$\$\$\$\$!
           :\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$^ ~\$\$\$\$\$\$\$\$\$\$\$\$  \$\$\$\$\$\$\$\$\$\$\$\$\$\$\$!
           :\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$~ ~\$\$\$\$\$\$\$\$\$\$\$\$  \$\$\$\$\$\$\$\$\$\$\$\$\$\$\$!
           :\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$~^^:.     ..:^~!\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$!
           ^\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$!
           :\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$~
            :\$\$\$\$\$\$\$\$\$\$\$\$\$:\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$~~\$\$\$\$\$\$\$\$\$\$\$\$~
              !\$\$\$\$\$\$\$\$\$\$. :\$\$..\$\$! :\$\$^ !\$!  ~\$\$\$\$\$\$\$\$\$\$.
               ^\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$!
                 \$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$:
                  ~\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$
                   "\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$\$"
                     \$\$\$\$\$~\$\$\$\$\$^\$\$\$\$\$~\$\$\$\$\$\$~\$\$\$\$:
                      \$\$^  .\$\$\$   \$\$\$:  ~\$\$^  .\$\$^
                      ..     :     :     :.     :
${reset}
"
}

detect_profile() {
  local DETECTED_PROFILE
  DETECTED_PROFILE=''
  local SHELLTYPE
  SHELLTYPE="$(basename "/$SHELL")"

  if [ "$SHELLTYPE" = "bash" ]; then
    if [ -f "$HOME/.bashrc" ]; then
      DETECTED_PROFILE="$HOME/.bashrc"
    elif [ -f "$HOME/.bash_profile" ]; then
      DETECTED_PROFILE="$HOME/.bash_profile"
    fi
  elif [ "$SHELLTYPE" = "zsh" ]; then
    DETECTED_PROFILE="$HOME/.zshrc"
  elif [ "$SHELLTYPE" = "fish" ]; then
    DETECTED_PROFILE="$HOME/.config/fish/conf.d/turso.fish"
  fi

  if [ -z "$DETECTED_PROFILE" ]; then
    if [ -f "$HOME/.profile" ]; then
      DETECTED_PROFILE="$HOME/.profile"
    elif [ -f "$HOME/.bashrc" ]; then
      DETECTED_PROFILE="$HOME/.bashrc"
    elif [ -f "$HOME/.bash_profile" ]; then
      DETECTED_PROFILE="$HOME/.bash_profile"
    elif [ -f "$HOME/.zshrc" ]; then
      DETECTED_PROFILE="$HOME/.zshrc"
    elif [ -d "$HOME/.config/fish" ]; then
      DETECTED_PROFILE="$HOME/.config/fish/conf.d/turso.fish"
    fi
  fi

  if [ ! -z "$DETECTED_PROFILE" ]; then
    echo "$DETECTED_PROFILE"
  fi
}

update_profile() {
   PROFILE_FILE=$(detect_profile)
   if [[ -n "$PROFILE_FILE" ]]; then
     if ! grep -q "\.turso" $PROFILE_FILE; then
        printf "\n${bright_blue}Updating profile ${reset}$PROFILE_FILE\n"
        printf "\n# Turso\nexport PATH=\"$INSTALL_DIRECTORY:\$PATH\"\n" >> $PROFILE_FILE
        printf "\nTurso will be available when you open a new terminal.\n"
        printf "If you want to make Turso available in this terminal, please run:\n"
        printf "\nsource $PROFILE_FILE\n"
     fi
   else
     printf "\n${bright_blue}Unable to detect profile file location. ${reset}Please add the following to your profile file:\n"
     printf "\nexport PATH=\"$INSTALL_DIRECTORY:\$PATH\"\n"
   fi
}

printf "\nWelcome to the Turso installer!\n"

print_logo
probe_arch
probe_os

URL_PREFIX="https://github.com/chiselstrike/homebrew-tap/releases/latest/download/"

TARGET="${OS}_$ARCH"

printf "${bright_blue}Downloading ${reset}$TARGET ...\n"

URL="$URL_PREFIX/homebrew-tap_$TARGET.tar.gz"

DOWNLOAD_FILE=$(mktemp -t turso.XXXXXXXXXX)

curl --progress-bar -L "$URL" -o "$DOWNLOAD_FILE"

INSTALL_DIRECTORY="$HOME/.turso"

printf "\n${bright_blue}Installing to ${reset}$INSTALL_DIRECTORY\n"

mkdir -p $INSTALL_DIRECTORY

tar -C $INSTALL_DIRECTORY -zxf $DOWNLOAD_FILE turso

rm -f $DOWNLOAD_FILE

update_profile

printf "\nTurso CLI installed!\n\n"
printf "If you are a new user, you can sign up with ${bright_blue}turso auth signup${reset}.\n\n"
printf "If you already have an account, please login with ${bright_blue}turso auth login${reset}.\n\n"

# DON'T RUN SIGNUP '$INSTALL_DIRECTORY/turso auth signup
