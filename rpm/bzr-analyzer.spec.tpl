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

#BuildRequires: java-1.8.0-openjdk-headless java-1.8.0-openjdk-devel java-1.8.0-openjdk maven chkconfig 

#%if 0%{?rhel} < 7
#Requires: bay-gcc61 gcc gcc-c++ clang clang-analyzer cppcheck oclint rats splint
#%else
#Requires: bay-gcc61 gcc gcc-c++ clang clang-analyzer cppcheck oclint rats splint frama-c pylint
#%endif

%description
The tool for simplification of using some code static analyzers such as cppcheck, oclint, rats etc

%if 0%{?rhel} >= 7
%package citool
Summary:  The tool for making SonarQube and bayzr integration
Requires: git systemd wget unzip shadow-utils cronie

%description citool
The tool for making SonarQube and bayzr integration
%endif

%prep
%setup -q

%build
export GOROOT=/usr/lib/golang/
export PATH=$PATH:$GOROOT/bin
export GOPATH=$(pwd)
export PATH=$PATH:$GOPATH/bin:$GOPATH/go-bindata/bin:$GOPATH/cisetup/bin
export GOPATH=$GOPATH:$GOPATH/go-bindata/:$GOPATH/cisetup/

/usr/bin/go build -o bin/bayzr main

#alternatives --install /usr/bin/java java /opt/jdk1.8.0_101/bin/java 2
#cd sonarqube
#mvn clean package
#cd ..

%if 0%{?rhel} >= 7
/usr/bin/go build -o bin/go-bindata github.com/jteeuwen/go-bindata/go-bindata

bin/go-bindata -o cisetup/src/data/data.go -pkg data cisetup/src/data cisetup/src/data/css cisetup/src/data/js cisetup/src/data/fonts cisetup/src/data/js/i18n
/usr/bin/go build -o bin/citool cisetup/src/main/main.go

%endif

%install
rm -rf $RPM_BUILD_ROOT
mkdir -p $RPM_BUILD_ROOT%{_bindir}
mkdir -p $RPM_BUILD_ROOT%{_sysconfdir}/bzr.d
mkdir -p $RPM_BUILD_ROOT%{_datadir}/doc/%{pkgname}/
mkdir -p $RPM_BUILD_ROOT%{_datarootdir}/bzr.java/

install -D -p -m 755 bin/bayzr %{buildroot}%{_bindir}
install -D -p -m 644 cfg/bzr.conf %{buildroot}%{_sysconfdir}
for f in cfg/*.conf; do
    fn=$(basename "$f")
    if [ "$fn" != "bzr.conf" ]; then
	install -D -p -m 644 $f %{buildroot}%{_sysconfdir}/bzr.d/$fn
    fi
done
%if 0%{?rhel} < 7
    rm -rf %{buildroot}%{_sysconfdir}/bzr.d/frama-c.conf
%endif
for f in cfg/*.tpl; do
    fn=$(basename "$f")
    if [ "$fn" != "bzr.conf.tpl" -a "$fn" != "checkerplugin.cfd.tpl" ]; then
	install -D -p -m 644 $f %{buildroot}%{_sysconfdir}/bzr.d/$fn
    fi
done
for f in xml/*.xml; do
    fn=$(basename "$f")
    install -D -p -m 644 $f %{buildroot}%{_sysconfdir}/bzr.d/$fn
done
install -D -p -m 644 rpm/gpl-3.0.txt %{buildroot}%{_datadir}/doc/%{pkgname}/LICENSE
install -D -p -m 644 rpm/LICENSE_GOCUI %{buildroot}%{_datadir}/doc/%{pkgname}/LICENSE_GOCUI
install -D -p -m 644 rpm/COPYRIGHT %{buildroot}%{_datadir}/doc/%{pkgname}/COPYRIGHT

install -D -p -m 644 sonarqube/jar/bayzr-plugin-0.0.1-rel1.jar %{buildroot}%{_datarootdir}/bzr.java/

%if 0%{?rhel} >= 7
#citool
mkdir -p $RPM_BUILD_ROOT%{_sbindir}
mkdir -p $RPM_BUILD_ROOT/mnt/chroot/
mkdir -p $RPM_BUILD_ROOT/usr/lib/systemd/system/
mkdir -p $RPM_BUILD_ROOT/usr/share/bzr.cron/
mkdir -p $RPM_BUILD_ROOT%{_sysconfdir}/cron.d
install -D -p -m 755 bin/citool %{buildroot}%{_sbindir}
install -D -p -m 600 cisetup/cfg/citool.ini %{buildroot}%{_sysconfdir}
install -D -p -m 644 cisetup/cfg/citool.service %{buildroot}/usr/lib/systemd/system/
install -D -p -m 755 cisetup/src/data/bayzr_clean_orph_environ.sh %{buildroot}/usr/share/bzr.cron/bayzr_clean_orph_environ.sh
install -D -p -m 644 cisetup/src/data/bayzr-citool %{buildroot}%{_sysconfdir}/cron.d/bayzr-citool
%endif

%if 0%{?rhel} >= 7
%pre citool
if [ $1 -eq 1 ] ; then
adduser checker
fi

%posttrans citool
/usr/bin/systemctl daemon-reload
setcap cap_sys_chroot+ep %{_sbindir}/citool
%endif

%clean
rm -rf $RPM_BUILD_ROOT

%if 0%{?rhel} >= 7
%files citool
%attr(600, checker, checker) %config(noreplace) %{_sysconfdir}/citool.ini
%defattr(-,root,root,-)
%dir /mnt/chroot/
%{_sbindir}/citool
/usr/lib/systemd/system/citool.service
/usr/share/bzr.cron/bayzr_clean_orph_environ.sh
%config(noreplace) %{_sysconfdir}/cron.d/bayzr-citool
%endif

%files
%defattr(-,root,root,-)
%doc %{_datadir}/doc/%{pkgname}/LICENSE
%doc %{_datadir}/doc/%{pkgname}/LICENSE_GOCUI
%doc %{_datadir}/doc/%{pkgname}/COPYRIGHT
%{_bindir}/bayzr
%config(noreplace) %{_sysconfdir}/bzr.conf
%config(noreplace) %{_sysconfdir}/bzr.d/*.conf
%{_sysconfdir}/bzr.d/*.tpl
%{_sysconfdir}/bzr.d/*.xml
%{_datarootdir}/bzr.java/*.jar

%changelog
