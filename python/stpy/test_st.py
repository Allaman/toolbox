import st


def test_steuerlast_gmbh():
    # Sample test cases
    test_cases = [
        # format: (gewinn, hebesatz, expected_output)
        (100000, 360, (15000, 825, 12600)),
        (50000, 200, (7500, 412.50, 3500)),
        (10000000, 460, (1500000, 82500, 1610000)),
    ]

    for gewinn, hebesatz, expected in test_cases:
        result = st.steuerlast_gmbh(gewinn, hebesatz)
        assert (
            result == expected
        ), f"Expected {expected} for gewinn {gewinn}, hebesatz {hebesatz} but got {result}"

    print("All tests passed!")


test_steuerlast_gmbh()
