This is an example of prometheus exporter that:
- list current running processes in linux host
- read /proc/[pid]/status
- export VM RSS of those processes to prometheus format at :8080/metrics


# Build
Run:

    make


# Run

    go run *.go
    curl localhost:8080/metrics

Should have output like below:

    # HELP process_memory_rss_bytes Size of memory resient set size of process read from /proc/[pid]/status
    # TYPE process_memory_rss_bytes gauge
    process_memory_rss_bytes{name="(sd-pam)",pid="758"} 2.92864e+06
    process_memory_rss_bytes{name="NetworkManager",pid="691"} 2.0226048e+07
    process_memory_rss_bytes{name="Xorg",pid="801"} 2.7267072e+07
    process_memory_rss_bytes{name="acpi_thermal_pm",pid="127"} 0
    process_memory_rss_bytes{name="at-spi-bus-laun",pid="834"} 5.967872e+06
    process_memory_rss_bytes{name="at-spi2-registr",pid="932"} 5.976064e+06
    process_memory_rss_bytes{name="auditd",pid="59"} 0
    process_memory_rss_bytes{name="avahi-daemon",pid="688"} 4.067328e+06
    process_memory_rss_bytes{name="avahi-daemon",pid="699"} 348160
    process_memory_rss_bytes{name="blockd",pid="70"} 0
    process_memory_rss_bytes{name="bluetoothd",pid="689"} 5.967872e+06

