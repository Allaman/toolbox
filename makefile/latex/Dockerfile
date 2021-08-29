FROM ubuntu:20.04

ARG USER

RUN ln -snf /usr/share/zoneinfo/Etc/UTC /etc/localtime \
    && echo "Etc/UTC" > /etc/timezone \
    DEBIAN_FRONTEND="noninteractive" \
    && apt-get -qq update < /dev/null > /dev/null \
    && apt-get -qq upgrade < /dev/null > /dev/null \
    && apt-get -qq install texlive-xetex texlive texlive-fonts-extra texlive-latex-extra < /dev/null > /dev/null \
    #&& apt-get -qq install texlive-latex-base texlive-latex-extra texlive-fonts-recommended xzdec < /dev/null > /dev/null \
    && rm -rf /var/lib/apt/lists/* \
    && useradd --create-home -u 1000 ${USER}

USER ${USER}
WORKDIR /home/${USER}

