def steuerlast_gmbh(gewinn, hebesatz):
    """Berechne die (vereinfachte) Steuerlast einer GmbH

    Args:
        gewinn (decimal): Der zu versteuernde Gewinn der GmbH
        hebesatz (decimal): Der Hebesatz der Gemeinde

    Returns:
        kst (decimal): Die Körperschaftsstuer
        soli (decimal): Der Solidaritätszuschlag
        gst (decimal): Die Gewerbesteuer
    """
    kst = gewinn * (15 / 100)
    soli = kst * (5.5 / 100)

    gst = gewinn * (3.5 / 100) * (hebesatz / 100)

    return (round(kst, 2), round(soli, 2), round(gst, 2))
