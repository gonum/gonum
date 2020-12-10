#! /usr/bin/env ruby
# Parse test manifest to create driver and area-specific test manifests

require 'getoptlong'
require 'csv'
require 'json'
require 'haml'
require 'fileutils'

class Manifest
  JSON_STATE = JSON::State.new(
    :indent       => "  ",
    :space        => " ",
    :space_before => "",
    :object_nl    => "\n",
    :array_nl     => "\n"
  )

  TITLE = {
    urgna2012: "RDF Graph Normalization (URGNA2012)",
    urdna2015: "RDF Dataset Normalization (URDNA2015)",
  }
  DESCRIPTION = {
    urgna2012: "Tests the 2012 version of RDF Graph Normalization.",
    urdna2015: "Tests the 2015 version of RDF Dataset Normalization."
  }

  Test = Struct.new(:id, :name, :comment, :approval, :action, :urgna2012, :urdna2015)

  attr_accessor :tests

  def initialize
    csv = CSV.new(File.open(File.expand_path("../manifest.csv", __FILE__)))

    columns = []
    csv.shift.each_with_index {|c, i| columns[i] = c.to_sym if c}

    @tests = csv.map do |line|
      entry = {}
      # Create entry as object indexed by symbolized column name
      line.each_with_index {|v, i| entry[columns[i]] = v ? v.gsub("\r", "\n").gsub("\\", "\\\\") : nil}

      urgna2012 = "#{entry[:test]}-urgna2012.nq" if entry[:urgna2012] == "TRUE"
      urdna2015 = "#{entry[:test]}-urdna2015.nq" if entry[:urdna2015] == "TRUE"
      Test.new(entry[:test], entry[:name], entry[:comment], entry[:approval],
               "#{entry[:test]}-in.nq", urgna2012, urdna2015)
    end
  end

  # Create files referenced in the manifest
  def create_files
    tests.each do |test|
      files = [test.action, test.urgna2012, test.urdna2015].compact
      files.compact.select {|f| !File.exist?(f)}.each do |f|
        File.open(f, "w") {|io| io.puts( f.end_with?('.json') ? "{}" : "")}
      end
    end
  end

  def test_class(test, variant)
    case variant.to_sym
    when :urgna2012 then "rdfn:Urgna2012EvalTest"
    when :urdna2015 then "rdfn:Urdna2015EvalTest"
    end
  end

  def to_jsonld(variant)
    context = ::JSON.parse %({
      "xsd": "http://www.w3.org/2001/XMLSchema#",
      "rdfs": "http://www.w3.org/2000/01/rdf-schema#",
      "mf": "http://www.w3.org/2001/sw/DataAccess/tests/test-manifest#",
      "mq": "http://www.w3.org/2001/sw/DataAccess/tests/test-query#",
      "rdfn": "http://json-ld.github.io/normalization/test-vocab#",
      "rdft": "http://www.w3.org/ns/rdftest#",
      "id": "@id",
      "type": "@type",
      "action": {"@id": "mf:action",  "@type": "@id"},
      "approval": {"@id": "rdft:approval", "@type": "@id"},
      "comment": "rdfs:comment",
      "entries": {"@id": "mf:entries", "@type": "@id", "@container": "@list"},
      "label": "rdfs:label",
      "name": "mf:name",
      "result": {"@id": "mf:result", "@type": "@id"}
    })

    manifest = {
      "@context" => context,
      "id" => "manifest-#{variant}",
      "type" => "mf:Manifest",
      "label" => TITLE[variant],
      "comment" => DESCRIPTION[variant],
      "entries" => []
    }

    tests.each do |test|
      next unless test.send(variant)

      manifest["entries"] << {
        "id" => "manifest-#{variant}##{test.id}",
        "type" => test_class(test, variant),
        "name" => test.name,
        "comment" => test.comment,
        "approval" => (test.approval ? "rdft:#{test.approval}" : "rdft:Proposed"),
        "action" => test.action,
        "result" => test.send(variant)
      }
    end

    manifest.to_json(JSON_STATE)
  end

  def to_html
    # Create vocab.html using vocab_template.haml and compacted vocabulary
    template = File.read("template.haml")
    manifests = TITLE.keys.inject({}) do |memo, v|
      memo["manifest-#{v}"] = ::JSON.load(File.read("manifest-#{v}.jsonld"))
      memo
    end

    Haml::Engine.new(template, :format => :html5).render(self,
      man: ::JSON.load(File.read("manifest.jsonld")),
      manifests: manifests
    )
  end

  def to_ttl(variant)
    output = []
    output << %(## RDF Dataset Normalization tests
## Distributed under both the W3C Test Suite License[1] and the W3C 3-
## clause BSD License[2]. To contribute to a W3C Test Suite, see the
## policies and contribution forms [3]
##
## 1. http://www.w3.org/Consortium/Legal/2008/04-testsuite-license
## 2. http://www.w3.org/Consortium/Legal/2008/03-bsd-license
## 3. http://www.w3.org/2004/10/27-testcases
##
## Test types
## * rdfn:Urgna2012EvalTest  - Normalization using URGNA2012
## * rdfn:Urdna2015EvalTest  - Normalization using URDNA2015

@prefix : <manifest-#{variant}#> .
@prefix rdf:  <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .
@prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> .
@prefix mf:   <http://www.w3.org/2001/sw/DataAccess/tests/test-manifest#> .
@prefix rdft: <http://www.w3.org/ns/rdftest#> .
@prefix rdfn: <http://json-ld.github.io/normalization/test-vocab#> .

<manifest-#{variant}>  a mf:Manifest ;
)
    output << %(  rdfs:label "#{TITLE[variant]}";)
    output << %(  rdfs:comment "#{DESCRIPTION[variant]}";)
    output << %(  mf:entries \()

    tests.select {|t| t.send(variant)}.map {|t| ":#{t.id}"}.each_slice(10) do |entries|
      output << %(    #{entries.join(' ')})
    end
    output << %(  \) .)

    tests.select {|t| t.send(variant)}.each do |test|
      output << "" # separator
      output << ":#{test.id} a #{test_class(test, variant)};"
      output << %(  mf:name "#{test.name}";)
      output << %(  rdfs:comment "#{test.comment}";) if test.comment
      output << %(  rdft:approval #{(test.approval ? "rdft:#{test.approval}" : "rdft:Proposed")};)
      output << %(  mf:action <#{test.action}>;)
      output << %(  mf:result <#{test.send(variant)}>;)
      output << %(  .)
    end
    output.join("\n")
  end
end

options = {
  output: $stdout
}

OPT_ARGS = [
  ["--format", "-f",  GetoptLong::REQUIRED_ARGUMENT,"Output format, default #{options[:format].inspect}"],
  ["--output", "-o",  GetoptLong::REQUIRED_ARGUMENT,"Output to the specified file path"],
  ["--quiet",         GetoptLong::NO_ARGUMENT,      "Supress most output other than progress indicators"],
  ["--touch",         GetoptLong::NO_ARGUMENT,      "Create referenced files and directories if missing"],
  ["--variant",       GetoptLong::REQUIRED_ARGUMENT,"Test variant, 'rdf' or 'json'"],
  ["--help", "-?",    GetoptLong::NO_ARGUMENT,      "This message"]
]
def usage
  STDERR.puts %{Usage: #{$0} [options] URL ...}
  width = OPT_ARGS.map do |o|
    l = o.first.length
    l += o[1].length + 2 if o[1].is_a?(String)
    l
  end.max
  OPT_ARGS.each do |o|
    s = "  %-*s  " % [width, (o[1].is_a?(String) ? "#{o[0,2].join(', ')}" : o[0])]
    s += o.last
    STDERR.puts s
  end
  exit(1)
end

opts = GetoptLong.new(*OPT_ARGS.map {|o| o[0..-2]})

opts.each do |opt, arg|
  case opt
  when '--format'       then options[:format] = arg.to_sym
  when '--output'       then options[:output] = File.open(arg, "w")
  when '--quiet'        then options[:quiet] = true
  when '--touch'        then options[:touch] = true
  when '--variant'      then options[:variant] = arg.to_sym
  when '--help'         then usage
  end
end

vocab = Manifest.new
vocab.create_files if options[:touch]
if options[:format] || options[:variant]
  case options[:format]
  when :jsonld  then options[:output].puts(vocab.to_jsonld(options[:variant]))
  when :ttl     then options[:output].puts(vocab.to_ttl(options[:variant]))
  when :html    then options[:output].puts(vocab.to_html)
  else  STDERR.puts "Unknown format #{options[:format].inspect}"
  end
else
  Manifest::TITLE.keys.each do |variant|
    %w(jsonld ttl).each do |format|
      File.open("manifest-#{variant}.#{format}", "w") do |output|
        output.puts(vocab.send("to_#{format}".to_sym, variant))
      end
    end
  end
  File.open("index.html", "w") do |output|
    output.puts(vocab.to_html)
  end
end
