# dur

dur is a CLI tool for calculation of durations.

![DUR](dur.png)

## Installation
with go installed run:
```bash
go install github.com/Oppodelldog/dur@latest
```

## Usage

```bash
# calculate with durations
> dur 12h - 1m + 60s
12h0m0s

# floating point durations
> dur 0,1666666666667h
10m0s

# default operation is addition
> dur 12h1m60s
12h2m0s

# units support
> dur 1h1m1s1ms1us1ns
1h1m1.001001001s

# parentheses support
> dur "2h-(1h30m)"
30m0s

# multiplication
> dur 5*4*12*8h
1920h0m0s

# division
> dur 40h/5
8h0m0s
```

### verbose output

```bash
# -p=h for printing calculations in a human readable form
> dur -p=h 12h - 1m + 60s
     12h0m0s -         1m0s =     11h59m0s
    11h59m0s +         1m0s =      12h0m0s
12h0m0s

# -p=n for  printing calculations in nanoseconds
> dur -p=n 12h - 1m + 60s
    43200000000000 -        60000000000 =     43140000000000
    43140000000000 +        60000000000 =     43200000000000
12h0m0s
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)