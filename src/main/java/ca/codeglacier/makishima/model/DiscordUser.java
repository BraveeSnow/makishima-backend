package ca.codeglacier.makishima.model;

import java.util.Date;

public record DiscordUser(
  long id,
  String accessToken,
  String refreshToken,
  Date expirationDate
) {}
