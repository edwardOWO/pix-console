# Some metadata required by an RPM package
Name: pix-console
Summary: pix config service
Version: 20240308
Release: 4
License: MIT

%description
PIX　Config Service


%install
# copy the executable to buildroot.
mkdir -p %{buildroot}/opt/pix-console
cp -rpdf /tmp/pix-console/* %{buildroot}/opt/pix-console
chmod 755 %{buildroot}/opt/pix-console/pix-console

# generate the systemd unit file to buildroot.
mkdir -p %{buildroot}/etc/systemd/system
cat <<EOF> %{buildroot}/etc/systemd/system/pix-console.service
[Unit]

[Install]
WantedBy=multi-user.target

[Service]
WorkingDirectory=/opt/pix-console
ExecStart=/opt/pix-console/pix-console
Restart=always
RestartSec=5
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=%n
EOF

%post
systemctl daemon-reload

%files
/opt/pix-console/*
/etc/systemd/system/pix-console.service

%config(noreplace) /opt/pix-console/config