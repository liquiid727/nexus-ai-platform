class ConsoleSmsProvider {
  async sendCode(phone, code) {
    console.log(`[SMS] to ${phone}: ${code}`);
    return true;
  }
}

function createProvider() {
  const provider = process.env.SMS_PROVIDER || 'console';
  switch (provider) {
    default:
      return new ConsoleSmsProvider();
  }
}

module.exports = { createProvider };

