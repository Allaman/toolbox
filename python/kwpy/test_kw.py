# https://www.finanz-tools.de/kalenderwochen-rechner
import kw


def test_to_first_day_in_iso_week():
    # Sample test cases
    test_cases = [
        # format: (input_year, input_week, expected_output)
        (2020, 1, (2019, 12, 30)),
        (2020, 52, (2020, 12, 21)),
        (2021, 1, (2021, 1, 4)),
        (2022, 52, (2022, 12, 26)),
        (2023, 1, (2023, 1, 2)),
        (2023, 42, (2023, 10, 16)),
        (2024, 1, (2024, 1, 1)),
    ]

    for year, week, expected in test_cases:
        result = kw.to_first_day_in_iso_week(year, week)
        assert (
            result == expected
        ), f"Expected {expected} for year {year}, week {week} but got {result}"

    print("All tests passed!")


def test_iso_week_number():
    # Sample test cases
    test_cases = [
        # format: (input_year, input_month, input_day, expected_output)
        (2023, 1, 1, (52)),
        (2023, 10, 21, (42)),
        (2023, 12, 31, (52)),
        (2022, 12, 31, (52)),
        (2020, 12, 31, (53)),
        (2020, 1, 1, (1)),
    ]

    for year, month, day, expected in test_cases:
        result = kw.iso_week_number(year, month, day)
        assert (
            result == expected
        ), f"Expected {expected} for year {year}, month {month}, day {day} but got {result}"

    print("All tests passed!")


test_to_first_day_in_iso_week()
test_iso_week_number()
