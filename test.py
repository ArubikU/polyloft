import time

class Cronometer:
    def __init__(self):
        self.start_time = None
        self.end_time = None
    
    def start(self):
        self.start_time = time.time()
    
    def stop(self):
        self.end_time = time.time()
    
    def elapsedMilliseconds(self):
        if self.start_time is None or self.end_time is None:
            return 0
        return (self.end_time - self.start_time) * 1000
    
    def elapsedFormatted(self):
        elapsed_ms = self.elapsedMilliseconds()
        # Format as HH:MM:SS.mmm
        hours, rem = divmod(elapsed_ms / 1000, 3600)
        minutes, seconds = divmod(rem, 60)
        return f"{int(hours):02}:{int(minutes):02}:{seconds:.3f}"

cron = Cronometer()
cron.start()
x = 1000000000 * 1000000000
print(x)
cron.stop()
print(cron.elapsedMilliseconds())
print(cron.elapsedFormatted())