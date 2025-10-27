"""
Portfolio Heat Calculator
Fast vectorized heat calculations using pandas for position sizing validation.

Usage:
    # From Excel VBA via =PY() formula
    result = heat_calculator.check_heat_caps(positions_df, 75, "Tech/Comm", 10000, 0.04, 0.015)

    # Standalone testing
    python heat_calculator.py
"""

import pandas as pd
import numpy as np
from typing import Dict, Optional, Tuple


def portfolio_heat_after(positions_df: pd.DataFrame, add_r: float) -> float:
    """
    Calculates total portfolio heat + proposed trade.

    Args:
        positions_df: DataFrame with columns [Ticker, Status, TotalOpenR]
        add_r: Additional R dollars from proposed trade

    Returns:
        Total heat in dollars

    Example:
        >>> positions = pd.DataFrame({
        ...     'Ticker': ['AAPL', 'MSFT'],
        ...     'Status': ['Open', 'Open'],
        ...     'TotalOpenR': [75.0, 50.0]
        ... })
        >>> heat = portfolio_heat_after(positions, 75.0)
        >>> print(heat)  # 200.0
    """
    if positions_df.empty:
        return add_r

    # Filter to open positions only
    open_positions = positions_df[positions_df['Status'] != 'Closed']

    if open_positions.empty:
        return add_r

    # Sum existing heat
    current_heat = open_positions['TotalOpenR'].sum()

    # Handle NaN values
    if pd.isna(current_heat):
        current_heat = 0.0

    return float(current_heat + add_r)


def bucket_heat_after(positions_df: pd.DataFrame, bucket: str, add_r: float) -> float:
    """
    Calculates bucket-specific heat + proposed trade.

    Args:
        positions_df: DataFrame with columns [Bucket, Status, TotalOpenR]
        bucket: Bucket name to filter by
        add_r: Additional R dollars from proposed trade

    Returns:
        Bucket heat in dollars

    Example:
        >>> positions = pd.DataFrame({
        ...     'Bucket': ['Tech/Comm', 'Tech/Comm', 'Healthcare'],
        ...     'Status': ['Open', 'Open', 'Open'],
        ...     'TotalOpenR': [75.0, 50.0, 60.0]
        ... })
        >>> heat = bucket_heat_after(positions, 'Tech/Comm', 75.0)
        >>> print(heat)  # 200.0 (75 + 50 + 75)
    """
    if positions_df.empty:
        return add_r

    # Filter to matching bucket and open positions
    bucket_positions = positions_df[
        (positions_df['Bucket'] == bucket) &
        (positions_df['Status'] != 'Closed')
    ]

    if bucket_positions.empty:
        return add_r

    # Sum existing bucket heat
    current_heat = bucket_positions['TotalOpenR'].sum()

    # Handle NaN values
    if pd.isna(current_heat):
        current_heat = 0.0

    return float(current_heat + add_r)


def check_heat_caps(
    positions_df: pd.DataFrame,
    add_r: float,
    bucket: str,
    equity: float,
    port_cap_pct: float,
    bucket_cap_pct: float
) -> Dict[str, any]:
    """
    Validates proposed trade against heat caps.

    Args:
        positions_df: DataFrame with columns [Ticker, Bucket, Status, TotalOpenR]
        add_r: Additional R dollars from proposed trade
        bucket: Bucket name for proposed trade
        equity: Account equity (E)
        port_cap_pct: Portfolio heat cap as decimal (e.g., 0.04 for 4%)
        bucket_cap_pct: Bucket heat cap as decimal (e.g., 0.015 for 1.5%)

    Returns:
        Dict with keys:
            - portfolio_ok: bool (True if under cap)
            - bucket_ok: bool (True if under cap)
            - portfolio_heat: float (total heat in $)
            - bucket_heat: float (bucket heat in $)
            - portfolio_cap: float (cap limit in $)
            - bucket_cap: float (cap limit in $)
            - portfolio_pct: float (heat as % of equity)
            - bucket_pct: float (bucket heat as % of equity)

    Example:
        >>> positions = pd.DataFrame({
        ...     'Ticker': ['AAPL', 'MSFT'],
        ...     'Bucket': ['Tech/Comm', 'Tech/Comm'],
        ...     'Status': ['Open', 'Open'],
        ...     'TotalOpenR': [150.0, 100.0]
        ... })
        >>> result = check_heat_caps(positions, 75, 'Tech/Comm', 10000, 0.04, 0.015)
        >>> print(result['portfolio_ok'])  # False (325/10000 = 3.25% < 4% cap)
        >>> print(result['bucket_ok'])  # False (325/10000 = 3.25% > 1.5% cap)
    """
    # Calculate heat levels
    port_heat = portfolio_heat_after(positions_df, add_r)
    buck_heat = bucket_heat_after(positions_df, bucket, add_r)

    # Calculate caps
    port_cap = equity * port_cap_pct
    buck_cap = equity * bucket_cap_pct

    # Calculate percentages
    port_pct = (port_heat / equity) * 100 if equity > 0 else 0
    buck_pct = (buck_heat / equity) * 100 if equity > 0 else 0

    return {
        'portfolio_ok': port_heat <= port_cap,
        'bucket_ok': buck_heat <= buck_cap,
        'portfolio_heat': float(port_heat),
        'bucket_heat': float(buck_heat),
        'portfolio_cap': float(port_cap),
        'bucket_cap': float(buck_cap),
        'portfolio_pct': float(port_pct),
        'bucket_pct': float(buck_pct),
        'room_portfolio': float(max(0, port_cap - (port_heat - add_r))),
        'room_bucket': float(max(0, buck_cap - (buck_heat - add_r)))
    }


def get_open_positions_summary(positions_df: pd.DataFrame) -> Dict[str, any]:
    """
    Returns summary statistics for open positions.

    Args:
        positions_df: DataFrame with positions data

    Returns:
        Dict with summary stats:
            - total_positions: int
            - total_heat: float
            - avg_heat_per_position: float
            - buckets: Dict[str, float] (heat by bucket)

    Example:
        >>> positions = pd.DataFrame({
        ...     'Ticker': ['AAPL', 'MSFT', 'JNJ'],
        ...     'Bucket': ['Tech/Comm', 'Tech/Comm', 'Healthcare'],
        ...     'Status': ['Open', 'Open', 'Open'],
        ...     'TotalOpenR': [75.0, 50.0, 60.0]
        ... })
        >>> summary = get_open_positions_summary(positions)
        >>> print(summary['total_positions'])  # 3
        >>> print(summary['total_heat'])  # 185.0
    """
    if positions_df.empty:
        return {
            'total_positions': 0,
            'total_heat': 0.0,
            'avg_heat_per_position': 0.0,
            'buckets': {}
        }

    # Filter to open positions
    open_positions = positions_df[positions_df['Status'] != 'Closed']

    if open_positions.empty:
        return {
            'total_positions': 0,
            'total_heat': 0.0,
            'avg_heat_per_position': 0.0,
            'buckets': {}
        }

    # Calculate summary stats
    total_positions = len(open_positions)
    total_heat = float(open_positions['TotalOpenR'].sum())
    avg_heat = total_heat / total_positions if total_positions > 0 else 0.0

    # Heat by bucket
    bucket_heat = open_positions.groupby('Bucket')['TotalOpenR'].sum().to_dict()
    bucket_heat = {k: float(v) for k, v in bucket_heat.items()}

    return {
        'total_positions': total_positions,
        'total_heat': total_heat,
        'avg_heat_per_position': avg_heat,
        'buckets': bucket_heat
    }


def calculate_max_position_size(
    positions_df: pd.DataFrame,
    bucket: str,
    equity: float,
    port_cap_pct: float,
    bucket_cap_pct: float
) -> Dict[str, float]:
    """
    Calculates maximum allowable position size given current heat.

    Args:
        positions_df: DataFrame with positions data
        bucket: Bucket for proposed trade
        equity: Account equity
        port_cap_pct: Portfolio cap percentage
        bucket_cap_pct: Bucket cap percentage

    Returns:
        Dict with:
            - max_r_portfolio: float (max R based on portfolio cap)
            - max_r_bucket: float (max R based on bucket cap)
            - max_r_combined: float (min of the two - actual limit)

    Example:
        >>> positions = pd.DataFrame({
        ...     'Ticker': ['AAPL'],
        ...     'Bucket': ['Tech/Comm'],
        ...     'Status': ['Open'],
        ...     'TotalOpenR': [200.0]
        ... })
        >>> max_size = calculate_max_position_size(positions, 'Tech/Comm', 10000, 0.04, 0.015)
        >>> print(max_size['max_r_combined'])  # Will show max allowable R
    """
    # Calculate current heat
    current_port_heat = portfolio_heat_after(positions_df, 0)
    current_buck_heat = bucket_heat_after(positions_df, bucket, 0)

    # Calculate caps
    port_cap = equity * port_cap_pct
    buck_cap = equity * bucket_cap_pct

    # Calculate room
    max_r_portfolio = max(0, port_cap - current_port_heat)
    max_r_bucket = max(0, buck_cap - current_buck_heat)

    return {
        'max_r_portfolio': float(max_r_portfolio),
        'max_r_bucket': float(max_r_bucket),
        'max_r_combined': float(min(max_r_portfolio, max_r_bucket))
    }


def validate_position_sizing(
    positions_df: pd.DataFrame,
    add_r: float,
    bucket: str,
    equity: float,
    port_cap_pct: float,
    bucket_cap_pct: float
) -> Tuple[bool, str]:
    """
    All-in-one validation function with detailed error message.

    Args:
        positions_df: DataFrame with positions data
        add_r: Proposed position size in R
        bucket: Bucket name
        equity: Account equity
        port_cap_pct: Portfolio cap percentage
        bucket_cap_pct: Bucket cap percentage

    Returns:
        Tuple of (is_valid: bool, message: str)

    Example:
        >>> positions = pd.DataFrame(...)
        >>> valid, msg = validate_position_sizing(positions, 75, 'Tech/Comm', 10000, 0.04, 0.015)
        >>> if not valid:
        ...     print(msg)  # "Portfolio heat would exceed cap: 425.0 > 400.0"
    """
    result = check_heat_caps(positions_df, add_r, bucket, equity, port_cap_pct, bucket_cap_pct)

    if not result['portfolio_ok']:
        return False, (
            f"Portfolio heat would exceed cap: "
            f"${result['portfolio_heat']:.2f} > ${result['portfolio_cap']:.2f} "
            f"({result['portfolio_pct']:.1f}% > {port_cap_pct * 100:.1f}%)"
        )

    if not result['bucket_ok']:
        return False, (
            f"Bucket '{bucket}' heat would exceed cap: "
            f"${result['bucket_heat']:.2f} > ${result['bucket_cap']:.2f} "
            f"({result['bucket_pct']:.1f}% > {bucket_cap_pct * 100:.1f}%)"
        )

    return True, "Position sizing OK"


def _test_heat_calculator():
    """
    Standalone test function for development/debugging.
    Run: python heat_calculator.py
    """
    print("=" * 70)
    print("Heat Calculator Test")
    print("=" * 70)

    # Create sample positions
    positions = pd.DataFrame({
        'Ticker': ['AAPL', 'MSFT', 'NVDA', 'JNJ'],
        'Bucket': ['Tech/Comm', 'Tech/Comm', 'Tech/Comm', 'Healthcare'],
        'Status': ['Open', 'Open', 'Open', 'Open'],
        'TotalOpenR': [75.0, 50.0, 100.0, 60.0]
    })

    print("\nSample Positions:")
    print(positions)

    # Test settings
    equity = 10000
    port_cap_pct = 0.04  # 4%
    bucket_cap_pct = 0.015  # 1.5%
    add_r = 75.0
    bucket = 'Tech/Comm'

    print(f"\nTest Trade:")
    print(f"  Add R: ${add_r}")
    print(f"  Bucket: {bucket}")
    print(f"  Equity: ${equity}")
    print(f"  Portfolio Cap: {port_cap_pct * 100}% (${equity * port_cap_pct})")
    print(f"  Bucket Cap: {bucket_cap_pct * 100}% (${equity * bucket_cap_pct})")

    # Run validation
    result = check_heat_caps(positions, add_r, bucket, equity, port_cap_pct, bucket_cap_pct)

    print("\nValidation Result:")
    print(f"  Portfolio OK: {'✅' if result['portfolio_ok'] else '❌'}")
    print(f"    Current Heat: ${result['portfolio_heat']:.2f} ({result['portfolio_pct']:.1f}%)")
    print(f"    Cap: ${result['portfolio_cap']:.2f}")
    print(f"    Room: ${result['room_portfolio']:.2f}")

    print(f"\n  Bucket OK: {'✅' if result['bucket_ok'] else '❌'}")
    print(f"    Current Heat: ${result['bucket_heat']:.2f} ({result['bucket_pct']:.1f}%)")
    print(f"    Cap: ${result['bucket_cap']:.2f}")
    print(f"    Room: ${result['room_bucket']:.2f}")

    # Test validation message
    valid, msg = validate_position_sizing(positions, add_r, bucket, equity, port_cap_pct, bucket_cap_pct)
    print(f"\n  Validation: {msg}")

    # Test summary
    summary = get_open_positions_summary(positions)
    print("\nPortfolio Summary:")
    print(f"  Total Positions: {summary['total_positions']}")
    print(f"  Total Heat: ${summary['total_heat']:.2f}")
    print(f"  Avg Heat/Position: ${summary['avg_heat_per_position']:.2f}")
    print(f"  Heat by Bucket:")
    for buck, heat in summary['buckets'].items():
        print(f"    {buck}: ${heat:.2f}")

    print("\n" + "=" * 70)


if __name__ == "__main__":
    _test_heat_calculator()
