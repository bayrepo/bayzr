%global bayzr_version VERRR
%global bayzr_release RELLL
%define pkgname bayzr

%undefine _missing_build_ids_terminate_build

Name:           %{pkgname}
Version:        %{bayzr_version}
Release:        %{bayzr_release}
Summary:        The tool for simplification of using some code static analyzers such as cppcheck, oclint, rats etc


License: GPLv3
Group: Development/Languages
Source0: %{name}-%{version}.tar.gz
URL: http://wiki.bayrepo.net
BuildRoot: %{_tmppath}/%{name}-%{version}-%{release}-root-%(%{__id_u} -n)

BuildRequires: golang >= 1.6.2
BuildRequires: golang-src >= 1.6.2
BuildRequires: golang-bin >= 1.6.2

Requires: bay-gcc61 gcc gcc-c++ clang clang-analyzer cppcheck oclint rats splint

%description
The tool for simplification of using some code static analyzers such as cppcheck, oclint, rats etc

%prep
%setup -q

%build
export GOROOT=/usr/lib/golang/
export PATH=$PATH:$GOROOT/bin
export GOPATH=$(pwd)
export PATH=$PATH:$GOPATH/bin

/usr/bin/go build -o bin/bayzr main

%install
rm -rf $RPM_BUILD_ROOT
mkdir -p $RPM_BUILD_ROOT%{_bindir}
mkdir -p $RPM_BUILD_ROOT%{_sysconfdir}/bzr.d
mkdir -p $RPM_BUILD_ROOT%{_datadir}/doc/%{pkgname}/

install -D -p -m 755 bin/bayzr %{buildroot}%{_bindir}
install -D -p -m 644 cfg/bzr.conf %{buildroot}%{_sysconfdir}
for f in cfg/*.conf; do
    fn=$(basename "$f")
    if [ "$fn" != "bzr.conf" ]; then
	install -D -p -m 644 $f %{buildroot}%{_sysconfdir}/bzr.d/$fn
    fi
done
for f in cfg/*.tpl; do
    fn=$(basename "$f")
    if [ "$fn" != "bzr.conf.tpl" -a "$fn" != "checkerplugin.cfd.tpl" ]; then
	install -D -p -m 644 $f %{buildroot}%{_sysconfdir}/bzr.d/$fn
    fi
done
install -D -p -m 644 rpm/gpl-3.0.txt %{buildroot}%{_datadir}/doc/%{pkgname}/LICENSE
install -D -p -m 644 rpm/COPYRIGHT %{buildroot}%{_datadir}/doc/%{pkgname}/COPYRIGHT

%clean
rm -rf $RPM_BUILD_ROOT

%files
%defattr(-,root,root,-)
%doc %{_datadir}/doc/%{pkgname}/LICENSE 
%doc %{_datadir}/doc/%{pkgname}/COPYRIGHT
%{_bindir}/bayzr
%config(noreplace) %{_sysconfdir}/bzr.conf
%config(noreplace) %{_sysconfdir}/bzr.d/*.conf
%{_sysconfdir}/bzr.d/*.tpl

%changelog
