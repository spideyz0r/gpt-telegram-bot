%global go_version 1.18.10
%global go_release go1.18.10

Name:           gpt-telegram-bot
Version:        0.0.2
Release:        1%{?dist}
Summary:        gpt-telegram-bot tool
License:        GPLv3
URL:            https://github.com/spideyz0r/gpt-telegram-bot
Source0:        %{url}/archive/refs/tags/v%{version}.tar.gz

BuildRequires:  golang >= %{go_version}
BuildRequires:  git

%description
gpt-telegram-bot is a telegram chat bot that connects to OpenAI's chat APIs

%global debug_package %{nil}

%prep
%autosetup -n %{name}-%{version}

%build
go build -v -o %{name} -ldflags=-linkmode=external

%check
go test

%install
install -Dpm 0755 %{name} %{buildroot}%{_bindir}/%{name}

%files
%{_bindir}/gpt-telegram-bot

%license LICENSE

%changelog
* Tue Jul 11 2023 spideyz0r <47341410+spideyz0r@users.noreply.github.com> 0.0.2-1
- Initial build

