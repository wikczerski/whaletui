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
        self.theme_path: Optional[str] = None

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

    def get_screen_content(self) -> str:
        """
        Get the current screen content.

        Returns:
            Current screen content as string
        """
        return self.get_output()

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
        Take a screenshot of the current screen.

        Args:
            filename: Name of the screenshot file
        """
        screenshots_dir = Path(__file__).parent / "screenshots"
        screenshots_dir.mkdir(exist_ok=True)

        screenshot_path = screenshots_dir / filename

        with open(screenshot_path, 'w', encoding='utf-8') as f:
            f.write(self.get_screen_content())

        self.logger.info(f"Screenshot saved to: {screenshot_path}")

    def navigate_to_view(self, view_name: str) -> bool:
        """
        Navigate to a specific view in the application.

        Args:
            view_name: Name of the view to navigate to

        Returns:
            True if navigation successful, False otherwise
        """
        view_commands = {
            'containers': 'c',
            'images': 'i',
            'volumes': 'v',
            'networks': 'n',
            'swarm': 's',
            'nodes': 'nodes',
            'services': 'services'
        }

        command = view_commands.get(view_name.lower())
        if not command:
            self.logger.error(f"Unknown view: {view_name}")
            return False

        self.send_text(command)
        time.sleep(1)  # Wait for view to load

        # Check if we're in the expected view
        return self.wait_for_screen(view_name, timeout=5)

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
        self.send_key('Enter')
        time.sleep(1)

        return True
