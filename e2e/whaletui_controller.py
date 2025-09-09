"""
WhaleTUI Controller for e2e testing using pexpect.
"""
import os
import time
import logging
from typing import Optional, List, Tuple
from pathlib import Path

# Use wexpect on Windows, pexpect on Unix
import os
if os.name == 'nt':  # Windows
    try:
        import wexpect as pexpect
    except ImportError:
        import pexpect
else:  # Unix/Linux
    import pexpect


class WhaleTUIController:
    """Controller for interacting with WhaleTUI during e2e tests."""

    def __init__(self, binary_path: str, timeout: int = 30):
        """
        Initialize the WhaleTUI controller.

        Args:
            binary_path: Path to the WhaleTUI executable
            timeout: Default timeout for operations in seconds
        """
        self.binary_path = binary_path
        self.timeout = timeout
        self.process: Optional[pexpect.spawn] = None
        self.logger = logging.getLogger(__name__)
        # Set default theme path to the e2e config theme
        self.theme_path: Optional[str] = "/app/whaletui/e2e/config/theme.yaml"

    def start(self, args: List[str] = None, env: dict = None) -> None:
        """
        Start the WhaleTUI application.

        Args:
            args: Command line arguments to pass to WhaleTUI
            env: Environment variables to set
        """
        if self.process:
            self.cleanup()

        cmd_args = [self.binary_path]
        if args:
            cmd_args.extend(args)

        # Add theme path if set
        if self.theme_path:
            cmd_args.extend(["--theme", self.theme_path])
            self.logger.info(f"Added theme path: {self.theme_path}")
        else:
            self.logger.warning("No theme path set!")

        # Set up environment
        test_env = os.environ.copy()
        if env:
            test_env.update(env)

        # Set terminal size for consistent testing
        test_env['LINES'] = '24'
        test_env['COLUMNS'] = '80'

        self.logger.info(f"Starting WhaleTUI with command: {' '.join(cmd_args)}")

        try:
            self.process = pexpect.spawn(
                cmd_args[0],
                args=cmd_args[1:],
                env=test_env,
                timeout=self.timeout
            )
            self.process.logfile_read = logging.getLogger('whaletui_output').handlers[0].stream if logging.getLogger('whaletui_output').handlers else None

            # Wait for the application to start
            time.sleep(2)

        except Exception as e:
            self.logger.error(f"Failed to start WhaleTUI: {e}")
            raise

    def send_key(self, key: str) -> None:
        """
        Send a key to the application.

        Args:
            key: Key to send (e.g., 'q', 'Enter', 'Tab', 'Esc')
        """
        if not self.process:
            raise RuntimeError("WhaleTUI process not started")

        key_mapping = {
            'Enter': '\r',
            'Tab': '\t',
            'Esc': '\x1b',
            'Space': ' ',
            'Up': '\x1b[A',
            'Down': '\x1b[B',
            'Right': '\x1b[C',
            'Left': '\x1b[D',
            'Home': '\x1b[H',
            'End': '\x1b[F',
            'PageUp': '\x1b[5~',
            'PageDown': '\x1b[6~',
            'F1': '\x1bOP',
            'F2': '\x1bOQ',
            'F3': '\x1bOR',
            'F4': '\x1bOS',
            'F5': '\x1b[15~',
            'F6': '\x1b[17~',
            'F7': '\x1b[18~',
            'F8': '\x1b[19~',
            'F9': '\x1b[20~',
            'F10': '\x1b[21~',
            'F11': '\x1b[23~',
            'F12': '\x1b[24~',
            'Ctrl+C': '\x03',
            'Ctrl+D': '\x04',
            'Ctrl+Z': '\x1a',
        }

        actual_key = key_mapping.get(key, key)
        self.process.send(actual_key)
        time.sleep(0.1)  # Small delay to allow processing

    def send_text(self, text: str) -> None:
        """
        Send text to the application.

        Args:
            text: Text to send
        """
        if not self.process:
            raise RuntimeError("WhaleTUI process not started")

        self.process.send(text)
        time.sleep(0.1)

    def expect(self, pattern: str, timeout: int = None) -> Tuple[int, str]:
        """
        Wait for a pattern to appear in the output.

        Args:
            pattern: Pattern to wait for (regex)
            timeout: Timeout in seconds (uses default if None)

        Returns:
            Tuple of (index, matched text)
        """
        if not self.process:
            raise RuntimeError("WhaleTUI process not started")

        timeout = timeout or self.timeout
        try:
            index = self.process.expect(pattern, timeout=timeout)
            return index, self.process.after
        except pexpect.TIMEOUT:
            self.logger.error(f"Timeout waiting for pattern: {pattern}")
            self.logger.error(f"Current output: {self.get_output()}")
            raise
        except pexpect.EOF:
            self.logger.error("Process ended unexpectedly")
            raise

    def expect_any(self, patterns: List[str], timeout: int = None) -> Tuple[int, str]:
        """
        Wait for any of the given patterns to appear.

        Args:
            patterns: List of patterns to wait for
            timeout: Timeout in seconds (uses default if None)

        Returns:
            Tuple of (index, matched text)
        """
        if not self.process:
            raise RuntimeError("WhaleTUI process not started")

        timeout = timeout or self.timeout
        try:
            index = self.process.expect(patterns, timeout=timeout)
            return index, self.process.after
        except pexpect.TIMEOUT:
            self.logger.error(f"Timeout waiting for any of patterns: {patterns}")
            self.logger.error(f"Current output: {self.get_output()}")
            raise
        except pexpect.EOF:
            self.logger.error("Process ended unexpectedly")
            raise

    def get_output(self) -> str:
        """
        Get the current output from the application.

        Returns:
            Current output as string
        """
        if not self.process:
            return ""

        before = self.process.before if isinstance(self.process.before, str) else str(self.process.before)
        after = self.process.after if isinstance(self.process.after, str) else str(self.process.after)
        return before + after

    def get_clean_screen_content(self) -> str:
        """
        Get the current screen content with proper line breaks and formatting.

        Returns:
            Clean screen content as string with proper line breaks
        """
        if not self.process:
            return ""

        # Get the raw output
        raw_output = self.get_output()

        # First, handle carriage returns by converting them to line breaks
        # This is important for TUI apps that use \r to position cursor
        processed_output = raw_output.replace('\r\n', '\n').replace('\r', '\n')

        # Split by line breaks and clean up
        lines = []
        for line in processed_output.split('\n'):
            # Remove ANSI escape sequences and control characters
            import re
            clean_line = re.sub(r'\x1b\[[0-9;]*m', '', line)  # Remove color codes
            clean_line = re.sub(r'\x1b\[[0-9;]*[A-Za-z]', '', clean_line)  # Remove other escape sequences
            clean_line = clean_line.replace('\x1b[0m', '').replace('\x1b[1m', '')
            lines.append(clean_line)

        # Filter out empty lines at the beginning and end
        while lines and not lines[0].strip():
            lines.pop(0)
        while lines and not lines[-1].strip():
            lines.pop()

        # Special handling for TUI content that uses box drawing characters
        # If we have a line with box drawing characters, try to split it properly
        processed_lines = []
        for line in lines:
            if '║' in line and ('╔' in line or '╗' in line or '╚' in line or '╝' in line):
                # This looks like a TUI table - try to split it into logical rows
                # The TUI content is all on one line because of carriage returns
                # We need to split it at logical boundaries to make it readable
                import re

                # Split the line into logical table rows
                # Look for patterns that indicate row boundaries
                table_rows = []

                # First, try to split by looking for container ID patterns
                # Pattern: ║ followed by spaces and a 12-character hex ID
                container_pattern = r'║\s+[a-f0-9]{12}'
                matches = list(re.finditer(container_pattern, line))

                if matches:
                    # Split at container ID positions, but keep the full row together
                    last_pos = 0
                    for i, match in enumerate(matches):
                        if match.start() > last_pos:
                            # Add content before this match
                            before_content = line[last_pos:match.start()].strip()
                            if before_content:
                                table_rows.append(before_content)

                        # Find the end of this row (next container ID or end of line)
                        next_start = matches[i + 1].start() if i + 1 < len(matches) else len(line)
                        row_content = line[match.start():next_start].strip()
                        if row_content:
                            table_rows.append(row_content)

                        last_pos = next_start

                    # Add any remaining content
                    if last_pos < len(line):
                        remaining_content = line[last_pos:].strip()
                        if remaining_content:
                            table_rows.append(remaining_content)
                else:
                    # Fallback: try to split by other patterns
                    # Look for header and footer patterns
                    header_pattern = r'╔[^║]*╗'
                    footer_pattern = r'╚[^║]*╝'

                    # Try to split by header first
                    header_match = re.search(header_pattern, line)
                    if header_match:
                        table_rows.append(header_match.group().strip())
                        remaining = line[header_match.end():].strip()
                        if remaining:
                            table_rows.append(remaining)
                    else:
                        # If no clear patterns, just add the line as is
                        table_rows.append(line)

                # Add all table rows as separate lines
                processed_lines.extend(table_rows)
            else:
                processed_lines.append(line)

        # Join lines and ensure proper formatting
        return '\n'.join(processed_lines)

    def get_screen_content(self) -> str:
        """
        Get the current screen content.

        Returns:
            Current screen content as string
        """
        return self.get_clean_screen_content()

    def wait_for_screen(self, expected_content: str, timeout: int = None) -> bool:
        """
        Wait for specific content to appear on screen.

        Args:
            expected_content: Content to wait for
            timeout: Timeout in seconds (uses default if None)

        Returns:
            True if content found, False otherwise
        """
        try:
            self.expect(expected_content, timeout)
            return True
        except pexpect.TIMEOUT:
            return False

    def is_running(self) -> bool:
        """
        Check if the application is still running.

        Returns:
            True if running, False otherwise
        """
        if not self.process:
            return False

        return self.process.isalive()

    def stop(self) -> None:
        """Stop the application gracefully."""
        if self.process and self.is_running():
            try:
                # Try to send Ctrl+C for graceful shutdown
                self.send_key('Ctrl+C')
                time.sleep(1)

                # If still running, force kill
                if self.is_running():
                    self.process.terminate()
                    time.sleep(0.5)

                    if self.is_running():
                        self.process.kill()

            except Exception as e:
                self.logger.warning(f"Error stopping process: {e}")

    def cleanup(self) -> None:
        """Clean up the controller and stop the application."""
        self.stop()
        if self.process:
            self.process.close()
            self.process = None

    def take_screenshot(self, filename: str) -> None:
        """
        Take a screenshot of the current screen as text content.

        Args:
            filename: Name of the screenshot file (will be saved as .txt)
        """
        screenshots_dir = Path("/app/test-data/screenshots")
        screenshots_dir.mkdir(exist_ok=True)

        # Ensure the filename has .txt extension for text content
        if not filename.endswith('.txt'):
            filename = filename.replace('.png', '.txt')

        screenshot_path = screenshots_dir / filename

        # Get screen content and format it properly line by line
        screen_content = self.get_screen_content()

        with open(screenshot_path, 'w', encoding='utf-8') as f:
            # Write each line separately for better readability
            lines = screen_content.split('\n')
            for line in lines:
                # Ensure each line ends with a newline
                f.write(line.rstrip() + '\n')

        self.logger.info(f"Text screenshot saved to: {screenshot_path}")

    def take_png_screenshot(self, filename: str) -> None:
        """
        Take an actual PNG screenshot of the terminal (requires additional setup).
        This is a placeholder for future PNG screenshot functionality.

        Args:
            filename: Name of the screenshot file (will be saved as .png)
        """
        screenshots_dir = Path("/app/test-data/screenshots")
        screenshots_dir.mkdir(exist_ok=True)

        # Ensure the filename has .png extension
        if not filename.endswith('.png'):
            filename = filename.replace('.txt', '.png')

        screenshot_path = screenshots_dir / filename

        # For now, save as text with PNG extension as fallback
        # TODO: Implement actual PNG screenshot using terminal screenshot libraries
        screen_content = self.get_screen_content()

        with open(screenshot_path, 'w', encoding='utf-8') as f:
            f.write("# PNG Screenshot (Text Fallback)\n")
            f.write("# This is a text representation of the screen content\n")
            f.write("# To get actual PNG screenshots, implement terminal screenshot library\n\n")
            for line in screen_content.split('\n'):
                f.write(line + '\n')

        self.logger.info(f"PNG screenshot (text fallback) saved to: {screenshot_path}")

    def navigate_to_view(self, view_name: str) -> bool:
        """
        Navigate to a specific view in the application using command mode.

        Args:
            view_name: Name of the view to navigate to

        Returns:
            True if navigation successful, False otherwise
        """
        view_commands = {
            'containers': 'containers',
            'images': 'images',
            'volumes': 'volumes',
            'networks': 'networks',
            'swarm': 'services',
            'nodes': 'nodes',
            'services': 'services'
        }

        command = view_commands.get(view_name.lower())
        if not command:
            self.logger.error(f"Unknown view: {view_name}")
            return False

        try:
            # Enter command mode with ':'
            self.send_key(':')
            time.sleep(0.5)  # Wait for command mode to be active

            # Send the navigation command
            self.send_text(command)
            time.sleep(0.5)

            # Press Enter to execute the command
            self.send_key('\r')
            time.sleep(1)  # Wait for view to load

            # Check if we're in the expected view
            return self.wait_for_screen(view_name, timeout=5)
        except Exception as e:
            self.logger.error(f"Error navigating to {view_name}: {e}")
            return False

    def search(self, search_term: str) -> bool:
        """
        Perform a search in the current view.

        Args:
            search_term: Term to search for

        Returns:
            True if search was successful, False otherwise
        """
        # Press '/' to open search
        self.send_key('/')
        time.sleep(0.5)

        # Type search term
        self.send_text(search_term)
        time.sleep(0.5)

        # Press Enter to execute search
        self.send_key('\r')
        time.sleep(1)

        return True
